package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
	"golang.org/x/tools/benchmark/parse"
)

func main() {
	unit := flag.String("unit", "NsPerOp", "Benchmark Unit: NsPerOp | AllocedBytesPerOp | AllocsPerOp")
	output := flag.String("output", "BenchmarkQuery", "Output file prefix (no extension)")
	flag.Parse()

	scan := bufio.NewScanner(os.Stdin)

	variants := []string{"Complex", "Mid", "Simple"}
	frameworks := []string{"sql", "gorm", "sqlt", "ent", "sqlc", "bun", "xorm"}

	data := map[string]map[string]float64{}
	for _, variant := range variants {
		data[variant] = map[string]float64{}
	}

	for scan.Scan() {
		line := scan.Text()

		if !strings.HasPrefix(line, "BenchmarkQuery/") {
			continue
		}

		b, err := parse.ParseLine(line)
		if err != nil {
			panic(err)
		}

		name := strings.TrimPrefix(b.Name, "BenchmarkQuery/")
		name = strings.TrimSuffix(name, "-12")

		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			continue
		}
		variant, framework := parts[0], parts[1]

		switch *unit {
		case "NsPerOp":
			data[variant][framework] = b.NsPerOp
		case "AllocedBytesPerOp":
			data[variant][framework] = float64(b.AllocedBytesPerOp)
		case "AllocsPerOp":
			data[variant][framework] = float64(b.AllocsPerOp)
		}
	}

	chart := charts.NewBar()
	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: *unit + " BenchmarkQuery",
		}),
		charts.WithAnimation(false),
		charts.WithInitializationOpts(opts.Initialization{
			BackgroundColor: "#FFFFFF",
		}),
	)

	chart.SetXAxis(frameworks)

	for _, variant := range variants {
		values := make([]opts.BarData, len(frameworks))
		for i, fw := range frameworks {
			values[i] = opts.BarData{Value: data[variant][fw]}
		}
		chart.AddSeries(variant, values)
	}

	filename := fmt.Sprintf("%s_%s.png", *output, *unit)
	if err := render.MakeChartSnapshot(chart.RenderContent(), filename); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Chart written to %s\n", filename)
}
