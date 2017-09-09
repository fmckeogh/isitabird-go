package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

var labels []string

func makeTensorFromImage(filename string) (*tf.Tensor, image.Image, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	r := bytes.NewReader(b)
	img, _, err := image.Decode(r)

	if err != nil {
		return nil, nil, err
	}

	tensor, err := tf.NewTensor(string(b))
	if err != nil {
		return nil, nil, err
	}
	graph, input, output, err := decodeJpegGraph()
	if err != nil {
		return nil, nil, err
	}
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, nil, err
	}
	return normalized[0], img, nil
}

func decodeJpegGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.ExpandDims(s,
		op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)),
		op.Const(s.SubScope("make_batch"), int32(0)))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func loadLabels() {
	file, err := Asset("models/ssd_mobilenet_v1_coco/labels.txt")
	if err != nil {
		log.Fatal(err)
	}

	s := string(file[:])

	reader := strings.NewReader(s)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		var newlabel string = regexp.MustCompile("[^a-z]+").ReplaceAllString(scanner.Text(), "")
		labels = append(labels, newlabel)
	}
}

/*
func executegraph(inputimage string) string {
	model, err := Asset("models/ssd_mobilenet_v1_coco/frozen_inference_graph.pb")
	if err != nil {
		log.Fatal(err)
	}
	labels, err := Asset("models/ssd_mobilenet_v1_coco/labels.txt")
	if err != nil {
		log.Fatal(err)
	}
	loadLabels(labels)

	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	tensor, i, err := makeTensorFromImage(inputimage)
	if err != nil {
		log.Fatal(err)
	}

	b := i.Bounds()
	img := image.NewRGBA(b)
	draw.Draw(img, b, i, b.Min, draw.Src)

	inputop := graph.Operation("image_tensor")

	o1 := graph.Operation("detection_scores")
	o2 := graph.Operation("detection_classes")
	o3 := graph.Operation("num_detections")

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			inputop.Output(0): tensor,
		},
		[]tf.Output{
			o1.Output(0),
			o2.Output(0),
			o3.Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}

	probabilities := output[1].Value().([][]float32)[0]
	classes := output[2].Value().([][]float32)[0]

	return probabilities
}
*/

func getinference(inputimage string) ([]float32, []float32) {
	model, err := Asset("models/ssd_mobilenet_v1_coco/frozen_inference_graph.pb")
	if err != nil {
		log.Fatal(err)
	}
	/*
		labels, err := Asset("models/ssd_mobilenet_v1_coco/labels.txt")
		if err != nil {
			log.Fatal(err)
		}
	*/
	loadLabels()

	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	tensor, i, err := makeTensorFromImage(inputimage)
	if err != nil {
		fmt.Println("Here")
		log.Fatal(err)
	}

	b := i.Bounds()
	img := image.NewRGBA(b)
	draw.Draw(img, b, i, b.Min, draw.Src)

	inputop := graph.Operation("image_tensor")

	o1 := graph.Operation("detection_boxes")
	o2 := graph.Operation("detection_scores")
	o3 := graph.Operation("detection_classes")
	o4 := graph.Operation("num_detections")

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			inputop.Output(0): tensor,
		},
		[]tf.Output{
			o1.Output(0),
			o2.Output(0),
			o3.Output(0),
			o4.Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}

	probabilities := (output[1].Value().([][]float32)[0])[0:5]
	classes := (output[2].Value().([][]float32)[0])[0:5]

	return probabilities, classes
}
