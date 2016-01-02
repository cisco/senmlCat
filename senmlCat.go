package main

import (
	"fmt"
	"encoding/json"
	"encoding/xml"
	"github.com/ugorji/go/codec"
	"flag"
	"os"
	"io/ioutil"
)

// TODO
// replace relative times with absolute.
// listen on HTTP , UDP , TCP
// wrtie to HTTP POST 
// write to influxdb
// write to Kafka
// support EXI ??
// read old v=1 senml and conver to v=3
// perfornance measurements and size measurements
// better CLI -if json -of xml -i tcp:1234 -o http 2345 

type SenMLRecord struct {
	XMLName *bool `json:"_,omitempty" xml:"senml"`
	
	BaseName string `json:"bn,omitempty"  xml:"bn,attr,omitempty"`
	BaseTime int `json:"bt,omitempty"  xml:"bt,attr,omitempty"`
	BaseUnit string `json:"bu,omitempty"  xml:"bu,attr,omitempty"`
	Version int `json:"ver,omitempty"  xml:"ver,attr,omitempty"`
	
	Name string `json:"n,omitempty"  xml:"n,attr,omitempty"`
	Unit string `json:"u,omitempty"  xml:"u,attr,omitempty"`
	Time int `json:"t,omitempty"  xml:"t,attr,omitempty"`
	UpdateTime int `json:"ut,omitempty"  xml:"ut,attr,omitempty"`

	Value *float64 `json:"v,omitempty"  xml:"v,attr,omitempty"`
	StringValue string `json:"vs,omitempty"  xml:"vs,attr,omitempty"`
	BoolValue *bool `json:"vb,omitempty"  xml:"vb,attr,omitempty"`
	
	Sum *float64 `json:"s,omitempty"  xml:"sum,,attr,omitempty"`
}  


type SenML struct {
	XMLName *bool `json:"_,omitempty" xml:"sensml"`
	Xmlns string `json:"_,omitempty" xml:"xmlns,attr"`
	
	Records []SenMLRecord  ` xml:"senml"`
}


func main() {
	var err error;

	doIndentPtr := flag.Bool("i", false, "indent output")
	
	doJSONPtr := flag.Bool("json", false, "output JSON formatted SENML ")
	doCborPtr := flag.Bool("cbor", false, "output CBOR formatted SENML ")
	doXMLPtr  := flag.Bool("xml",  false, "output XML formatted SENML ")
	
	doIJSONPtr := flag.Bool("ijson", false, "input JSON formatted SENML ")
	doIXMLPtr := flag.Bool("ixml", false, "input XML formatted SENML ")
	doICBORPtr := flag.Bool("icbor", false, "input CBOR formatted SENML ")
	
	flag.Parse()

	// load the input  
	msg, err := ioutil.ReadFile( flag.Arg(0) )
	if err != nil {
		fmt.Printf("error reading JSON XML %v\n",err)
		os.Exit( 1 )
	}	
    //fmt.Print(string(msg))
	
	var s SenML
	s.XMLName = nil
	s.Xmlns = "urn:ietf:params:xml:ns:senml"
	
	// parse the input JSON
	if ( *doIJSONPtr ) {
		err = json.Unmarshal(msg, &s.Records )
		if err != nil {
			fmt.Printf("error parsing JSON XML %v\n",err)
			os.Exit( 1 )
		}
	}
	
	// parse the input XML
	if ( *doIXMLPtr ) {
		err = xml.Unmarshal(msg, &s)
		if err != nil {
			fmt.Printf("error parsing JSON XML %v\n",err)
			os.Exit( 1 )
		}
	}

	// parse the input CBOR
	if ( *doICBORPtr ) {
		var cborHandle codec.Handle = new( codec.CborHandle )
		var decoder *codec.Decoder = codec.NewDecoderBytes( msg, cborHandle )
		err = decoder.Decode( &s.Records )
			if err != nil {
			fmt.Printf("error parsing JSON XML %v\n",err)
			os.Exit( 1 )
		}
	}
	
	// ouput JSON version 
	if ( *doJSONPtr ) {
		var d []byte;
		if ( *doIndentPtr ) {
			d,err = json.MarshalIndent( s.Records, "", "  " )
		} else {
			d,err = json.Marshal( s.Records )
		}
		if err != nil {
			fmt.Printf("error encoding json %v\n",err)
			os.Exit( 1 )
		}
		fmt.Printf("%s\n", d)
	}

	// output a XML version 
	if ( *doXMLPtr ) {
		var d []byte;
		if ( *doIndentPtr ) {
			d,err = xml.MarshalIndent( s, "", "  " )
		} else {
			d,err = xml.Marshal( s )
		}
		if err != nil {
			fmt.Printf("error encoding xml %v\n",err);	
		}
		fmt.Printf("%s\n", d)
	}

	// output a CBOR version 
	if ( *doCborPtr ) {
		var d []byte 
		var cborHandle codec.Handle = new(codec.CborHandle)
		var encoder *codec.Encoder = codec.NewEncoderBytes( &d, cborHandle)
		err = encoder.Encode( s.Records )
		if err != nil {
			fmt.Printf("error encoding json %v\n",err)
			os.Exit( 1 )
		}
		fmt.Printf("%s\n", d)
	}
}

