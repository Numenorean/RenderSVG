package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/valyala/fasthttp"
)

type svg2pngRequests struct {
	SvgData string `json:"svg"`
}

var n = 0
const port = "26497"

func svg2png(svgData string, w, h, R, G, B, A int, addBackground bool) string {
	var in []byte
	in, err := base64.StdEncoding.DecodeString(svgData)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(in)

	icon, _ := oksvg.ReadIconStream(reader)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	if addBackground {
		draw.Draw(rgba, rgba.Bounds(), &image.Uniform{
			color.RGBA{
				uint8(R),
				uint8(G),
				uint8(B),
				uint8(A),
			},
		}, image.Point{}, draw.Src)
	}

	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	var out []byte
	writer := bytes.NewBuffer(out)
	png.Encode(writer, rgba)
	finalData, err := ioutil.ReadAll(writer)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(finalData)
}

func svg2pngHandler(ctx *fasthttp.RequestCtx) {
	var data svg2pngRequests

	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		fmt.Fprintf(ctx, "{\"error\":\"%q\"}", err)
		return
	}

	query := ctx.QueryArgs()
	width, _ := strconv.Atoi(string(query.Peek("width")))
	height, _ := strconv.Atoi(string(query.Peek("height")))
	r, _ := strconv.Atoi(string(query.Peek("r")))
	g, _ := strconv.Atoi(string(query.Peek("g")))
	b, _ := strconv.Atoi(string(query.Peek("b")))
	a, _ := strconv.Atoi(string(query.Peek("a")))
	addBackground, _ := strconv.ParseBool(string(query.Peek("addBackground")))

	png := svg2png(data.SvgData, width, height, r, g, b, a, addBackground)
	fmt.Fprintf(ctx, "{\"data\":\"%s\"}", png)
	n++
}

func main() {
	r := router.New()
	r.POST("/svg2png", svg2pngHandler)
	fmt.Println("Server started...")
	go func() {
		for true {
			fmt.Printf("\rRequest #%d", n)
			time.Sleep(time.Second * 1)
		}
	}()
	fasthttp.ListenAndServe(":"+port, r.Handler)
}
