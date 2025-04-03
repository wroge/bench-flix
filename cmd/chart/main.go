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
	benchmark := flag.String("benchmark", "BenchmarkQuery", "Benchmark Name")

	flag.Parse()

	scan := bufio.NewScanner(os.Stdin)

	frameworks := []string{"sql", "gorm", "sqlt", "ent", "sqlc", "bun", "xorm"}
	variants := []string{}

	data := map[string]map[string]float64{}

	for scan.Scan() {
		line := scan.Text()

		if !strings.HasPrefix(line, "Benchmark"+*benchmark+"/") {
			continue
		}

		b, err := parse.ParseLine(line)
		if err != nil {
			panic(err)
		}

		name := strings.TrimPrefix(b.Name, "Benchmark"+*benchmark+"/")
		name = strings.TrimSuffix(name, "-12")

		parts := strings.SplitN(name, "_", 2)

		var (
			variant, framework string
		)

		if len(parts) != 2 {
			framework = parts[0]
		} else {
			variant = parts[0]
			framework = parts[1]
		}

		if data[variant] == nil {
			variants = append(variants, variant)
			data[variant] = map[string]float64{}
		}

		switch *unit {
		case "NsPerOp":
			data[variant][framework] = b.NsPerOp
		case "AllocedBytesPerOp":
			data[variant][framework] = float64(b.AllocedBytesPerOp)
		case "AllocsPerOp":
			data[variant][framework] = float64(b.AllocsPerOp)
		}
	}

	fmt.Println(data)

	chart := charts.NewBar()
	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: *unit + " " + *benchmark,
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

	filename := fmt.Sprintf("%s_%s.png", *benchmark, *unit)
	if err := render.MakeChartSnapshot(chart.RenderContent(), filename); err != nil {
		panic(err)
	}

	fmt.Printf("Chart written to %s\n", filename)
}
