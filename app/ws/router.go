package ws

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
	"github.com/huobirdcenter/huobi_golang/pkg/response/common"
	"github.com/richard-xtek/go-grpc-micro-kit/log"
	"go.uber.org/zap"

	"go-prj-skeleton/app/auth"
	"go-prj-skeleton/app/config"
	"go-prj-skeleton/app/ws/consumer"
	"go-prj-skeleton/app/ws/handler"
	"go-prj-skeleton/app/ws/hub"
)

// Router ...
func Router(logger log.Factory, h *hub.Hub) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)

	rawStream := handler.NewRawStreamHandler(logger, h)
	router.HandleFunc("/ws/{stream_name}", rawStream.Handle)

	combinedStream := handler.NewCombinedStreamHandler(logger, h)
	router.HandleFunc("/stream", combinedStream.Handle)

	return router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

type srv struct {
	logger           log.Factory
	port             string
	huobiConfig      config.HuobiAPI
	kafkaConfig      *sarama.Config
	kafkaBrokers     []string
	signingSecretKey []byte
	server           *http.Server
}

// New ...
func New(port string, logger log.Factory, huobiConfig config.HuobiAPI, kafkaConfig *sarama.Config, brokers []string, signingSecretKey []byte) *srv {
	return &srv{
		port:             port,
		logger:           logger,
		huobiConfig:      huobiConfig,
		kafkaConfig:      kafkaConfig,
		kafkaBrokers:     brokers,
		signingSecretKey: signingSecretKey,
	}
}

func (s *srv) Start() error {
	h := hub.New(s.logger, auth.NewTokenValidator(s.signingSecretKey))

	// Don't remove: This used to test local
	// for i, stream := range []string{"btcusdt@depth", "htusdt@kline_1min", "htusdt@kline_3min", "htusdt@kline_1day", "htusdt@orderUpdate"} {
	// 	payload := struct {
	// 		Message string `json:"message"`
	// 	}{fmt.Sprintf("Message %dth sent", i)}

	// 	if stream == "htusdt@orderUpdate" {
	// 		stream = "@userData"
	// 		payload.Message += " for user 12345"
	// 		go h.StartMessage(payload, stream, "12345")
	// 		pay1 := payload
	// 		pay1.Message += " for user 123456"
	// 		go h.StartMessage(pay1, stream, "123456")
	// 	} else {
	// 		go h.StartMessage(payload, stream, "")
	// 	}
	// }

	commonClient := new(client.CommonClient).Init(s.huobiConfig.Host)

	symbols, err := commonClient.GetSymbols()
	if err != nil {
		return err
	}

	currencies := strings.Split(s.huobiConfig.BaseCurrency, ",")
	if len(currencies) == 0 {
		return fmt.Errorf("not base currency input")
	}

	symbols = s.filterByBaseCurrencies(symbols, currencies)

	cfg := &consumer.Config{
		Brokers: s.kafkaBrokers,
		Conf:    s.kafkaConfig,
	}
	consumer.StartConsumerers(context.Background(),
		consumer.NewCanndelstickConsumer(cfg, h, symbols, s.logger),
		consumer.NewMarketDepthConsumer(cfg, h, symbols, s.logger),
		consumer.NewMarketDetailConsumer(cfg, h, symbols, s.logger),
		consumer.NewUserEventConsumer(cfg, h, s.logger),
		consumer.NewTradeDetailConsumer(cfg, h, symbols, s.logger))

	router := Router(s.logger, h)

	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		s.logger.Bg().Error("Websocket server fail to start", zap.Error(err))
		return fmt.Errorf("start ws fail: %v", err)
	}

	server.RegisterOnShutdown(s.gracefullyShutdown)

	s.server = server

	return nil
}

func (s *srv) gracefullyShutdown() {
	// will call close current opened websocket connection ...
}

func (s *srv) Stop() error {
	if err := s.server.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func (s *srv) filterByBaseCurrencies(symbols []common.Symbol, currencies []string) []common.Symbol {
	out := []common.Symbol{}

	for _, cur := range currencies {
		filterSyms := filterByBaseCurrency(symbols, cur)
		if len(filterSyms) == 0 {
			err := fmt.Errorf("invalid base currency: %s", cur)
			s.logger.Bg().Error(fmt.Sprintf("don't have any symbols for base currency: %s", cur), zap.Error(err))
			panic(err)
		}

		out = append(out, filterSyms...)
	}

	return out
}

func filterByBaseCurrency(symbols []common.Symbol, currency string) []common.Symbol {
	out := []common.Symbol{}

	for _, sym := range symbols {
		if sym.BaseCurrency == currency {
			out = append(out, sym)
		}
	}

	return out
}
