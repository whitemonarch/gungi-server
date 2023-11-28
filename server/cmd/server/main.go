package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/spf13/viper"

	"github.com/whitemonarch/gungi-server/server/internal/api"
	"github.com/whitemonarch/gungi-server/server/internal/auth"
	"github.com/whitemonarch/gungi-server/server/internal/ws"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DB_SOURCE           string `mapstructure:"DB_SOURCE"`
	SUPABASE_JWT_SECRET string `mapstructure:"SUPABASE_JWT_SECRET"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {
	config, err := LoadConfig("./")
	if err != nil {
		log.Fatalln("Cannot load config", err)
	}

	db := api.DBConn{}
	maxRetries := 5
	sleepDuration := 2 * time.Second
	for i := 1; i <= maxRetries; i++ {
		err = db.PostgresConnect(config.DB_SOURCE)
		if err == nil {
			break
		}
		if i < maxRetries {
			log.Println("Connection failed, retrying...")
			time.Sleep(sleepDuration)
			sleepDuration *= 2
		}
	}
	if err != nil {
		log.Fatalln("Failed to establish a database connection: ", err)
	}
	defer db.Conn.Close()

	e := echo.New()
	m := melody.New()
	m.Config.MaxMessageSize = 1024

	// e.Use(middleware.VerifySupabaseTokenMiddleware)
	verify := e.Group("", VerifySupabaseTokenMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	e.GET("/game/:id", db.GetGameWithUndoRoute)
	verify.GET("/getongoinggamelist", db.GetOngoingGameList)

	e.GET("/ws", ws.WS(m, &db))

	verify.GET("/user/onboarding", db.GetOnboarding)
	verify.PUT("/user/onboarding", db.PutOnboarding)
	verify.PUT("/user/changename", db.ChangeUsername)

	e.Logger.Fatal(e.Start("localhost:5080"))
}

func VerifySupabaseTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.AuthenticateSupabaseToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}

		c.Set("sub", claims["sub"])

		return next(c)
	}
}
