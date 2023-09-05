package multichecker

import (
	"github.com/BurntSushi/toml"
	"github.com/jbakhtin/rtagent/pkg/errcheck"
	"github.com/jbakhtin/rtagent/pkg/osexitcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"
	"os"
)

// Статический анализатор.
// Включает в себя:
// - все анализаторы пакета analysis.
// - все анализаторы класса SA пакета staticcheck.
// - самописный анализатор на обнаружение необработанных ошибок.
// - самописный анализатор на обнаружение вызовов os.Exit в функции main.

// Для запуска необходимо:
//
//  1. Скомпилировать приложение
//     go build -o bin/multichecker cmd/staticlint/main.go
//
//  2. Запустить скомпилированный бинарник из консоли для указанной директории
//     ./bin/multichecker ./...
//
//   Для игнорирования проверки авто-сгенерированных файлов с тестами.
//   ./bin/multichecker -test=false ./...

// StaticCheckPath - путь к фалу конфигурации для staticcheck в фомате toml.
var StaticCheckPath = "staticcheck.toml"

// Config - структура файла конфигурации.
type Config struct {
	Checks []string
}

func Start() error {
	data, err := os.ReadFile(StaticCheckPath)
	if err != nil {
		return err
	}

	var cfg Config
	if err = toml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		tests.Analyzer,
		unreachable.Analyzer,

		osexitcheck.Analyzer,

		errcheck.Analyzer,
		fieldalignment.Analyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.Checks {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)

	return nil
}
