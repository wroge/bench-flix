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
	variants := flag.String("variants", "", "Benchmark Variants")
	frameworks := flag.String("frameworks", "sql,gorm,sqlt,ent,sqlc,bun", "Frameworks")

	flag.Parse()

	scan := bufio.NewScanner(os.Stdin)

	frameworksSlice := strings.Split(*frameworks, ",")

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

	chart.SetXAxis(frameworksSlice)

	variantSlice := strings.Split(*variants, ",")
	if len(variantSlice) == 0 {
		variantSlice = []string{""}
	}

	for _, variant := range variantSlice {
		values := make([]opts.BarData, len(frameworksSlice))
		for i, fw := range frameworksSlice {
			values[i] = opts.BarData{Value: data[variant][fw]}
		}

		chart.AddSeries(variant, values)
	}

	filename := fmt.Sprintf("%s_%s", *benchmark, *unit)

	if *variants != "" {
		filename += "_" + strings.ReplaceAll(*variants, ",", "")
	}

	if err := render.MakeChartSnapshot(chart.RenderContent(), "charts/"+filename+".png"); err != nil {
		panic(err)
	}

	fmt.Printf("Chart written to %s\n", filename)
}
