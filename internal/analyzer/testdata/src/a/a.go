package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestLogs() {
	logger := zap.NewNop()


	// - - - Simple init cases - - -


	// rule 1: 1st letter must be lowercase
	slog.Info("Starting server on port 8080") // want "log message must start with lowercase letter"
	slog.Error("Failed to connect to database")	// want "log message must start with lowercase letter"

	slog.Info("starting server on port 8080") // OK
	slog.Error("failed to connect to database")	// OK

	// rule 2: only english language
	slog.Info("запуск сервера") 					// want "log message must contain only english letters.*"
	slog.Error("ошибка подключения к базе данных")	// want "log message must contain only english letters.*"

	slog.Info("starting server")				// OK
	slog.Error("failed to connect to database")	// OK

	// rule 3: special characters, emoji
	slog.Info("server started! 🚀")					// want "log message must contain only english letters.*"
	slog.Error("connection failed!!!")				// want "log message should not end with punctuation marks"
	slog.Warn("warning: something went wrong...")	// want "log message should not end with punctuation marks"

	slog.Info("server started")			// OK
	slog.Error("connection failed")		// OK
	slog.Warn("something went wrong")	// OK

	// rule 4: sensitive data
	password := "pass123"
	slog.Info("user password:" + password)	// want "log message contains potentially sensitive data.*"
	apiKey := "key123"
	slog.Debug("api_key=" + apiKey)			// want "log message contains potentially sensitive data.*"
	token := "token123"
	slog.Info("token:" + token)				// want "log message contains potentially sensitive data.*"

	slog.Info("user authenticated successfully")	// OK
	slog.Debug("api request completed")				// OK
	slog.Info("token validated")					// OK


	// - - - Complex cases - - -


	// rule 1: 1st letter must be lowercase
	slog.Info("a")	// OK
	slog.Info("A")	// want "log message must start with lowercase letter"

	// rule 2: only english language
	slog.Info("test word")			// OK
	slog.Info("тестовое слово")		// want "log message must contain only english letters.*"
	slog.Info("测试词")				 // want "log message must contain only english letters.*"
	slog.Info("テスト単語")			  // want "log message must contain only english letters.*"
	slog.Info("테스트 단어")		   // want "log message must contain only english letters.*"
	slog.Info("słowo testowe")		// want "log message must contain only english letters.*"
	slog.Info("testovací slovo")	// want "log message must contain only english letters.*"
	slog.Info("δοκιμαστική λέξη")	// want "log message must contain only english letters.*"
	slog.Info("परीक्षण शब्द")			// want "log message must contain only english letters.*"
	slog.Info("คำทดสอบ")			// want "log message must contain only english letters.*"
	slog.Info("từ kiểm tra")		// want "log message must contain only english letters.*"

	// rule 3: special characters, emoji
	slog.Info("server started") // OK

	slog.Info("server started ✈️")						// want "log message must contain only english letters.*"
	slog.Info("ser‍ver started")						// want "log message must contain only english letters.*"
	slog.Info("server​started")							// want "log message must contain only english letters.*"
	slog.Info("server started ‮error")			// want "log message must contain only english letters.*"
	slog.Info("server started ‎‏error")	// want "log message must contain only english letters.*"
	slog.Info("se̷r̷v̷e̷r̷ started")						// want "log message must contain only english letters.*"
	slog.Error("s̴̩͝é̷̗r̴͉̓v̶͕̕e̴̠͐r̶̯̾ ̸̦͊f̴̛͖a̴͚͘i̵͉̍l̵̙͘e̷̳͝d̷͔̍")							// want "log message must contain only english letters.*"
	slog.Warn("server ∑ started ∆ connection ∞")		// want "log message must contain only english letters.*"
	slog.Info("ｓｅｒｖｅｒ　ｓｔａｒｔｅｄ")				// want "log message must contain only english letters.*"
	slog.Error("<<<<< server crashed >>>>>")			// want "log message must contain only english letters.*"
	slog.Info("server started 👨‍💻")						// want "log message must contain only english letters.*"
	slog.Info("server started 👍🏽")						// want "log message must contain only english letters.*"
	slog.Warn("server started ")						// want "log message must contain only english letters.*"
	slog.Info("server started ")						// want "log message must contain only english letters.*"
	slog.Info("ѕerver started")							// want "log message must contain only english letters.*"
	slog.Info("server sтarted")							// want "log message must contain only english letters.*"
	slog.Info("😀😃😁😂🤣😊😇🙂😌😍🥰😘🤓😎🥸🤩🥳🥺😢😭😤😠😡🤬🤯😳🥵🥶😱😨😰🤐🥴🤢🤮🤧😷🤒🤕🤑🤠😈👿👹👺")	// want "log message must contain only english letters.*"

	// rule 4: sensitive data
	slog.Info("user (id=123) logged in, path: /api/v1")	// OK
	logger.Info("100% processing done [status: ok]")	// OK
	slog.Info("user password: " + password)				// want "log message contains potentially sensitive data.*"
}
