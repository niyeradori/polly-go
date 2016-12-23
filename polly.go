package polly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bmizerany/aws4"
)

//http://docs.aws.amazon.com/polly/latest/dg/API_DescribeVoices.html
// SpeechResponse is the resource representing response from SynthesizeSpeech action.
type SpeechResponse struct {
	Audio       []byte
	RequestID   string
	ContentType string
}

//http://docs.aws.amazon.com/polly/latest/dg/API_DescribeVoices.html
// Voice is the structure element of ListResponse structure.
type Voice struct {
	Gender       string
	Id           string
	LanguageCode string
	LanguageName string
	Name         string
}

// VoiceResponse is the resource representing response from DescribeVoices action.
type VoicesResponse struct {
	NextToken string
	Voices    []Voice
}

// SpeechOptions is the set of parameters that can be used on the SynthesizeSpeech action.
// For more details see http://docs.aws.amazon.com/polly/latest/dg/API_SynthesizeSpeech.html
type SpeechOptions struct {
	OutputFormat string
	Text         string
	VoiceId      string
}

type VoiceOptions struct {
	LanguageCode string
	NextToken    string
}

// Polly is used to invoke API calls
type Polly struct {
	AccessKey string
	SecretKey string
}

// Polly Speech Cloud URL
const PollyAPI = "https://polly.us-west-2.amazonaws.com"
const synthesizeSpeechAPI = PollyAPI + "/v1/speech"
const describeVoicesAPI = PollyAPI + "/v1/voices"

// New returns a new Polly client.
func New(accessKey string, secretKey string) *Polly {
	return &Polly{AccessKey: accessKey, SecretKey: secretKey}
}

// NewSpeechOptions is the set of default parameters that can be used the SYntesizeSpeech action.
// For more details see http://docs.aws.amazon.com/polly/latest/dg/API_SynthesizeSpeech.html
func NewSpeechOptions(data string) SpeechOptions {
	return SpeechOptions{
		OutputFormat: "mp3",
		Text:         data,
		VoiceId:      "Joanna"}
}

// SynthesizeSpeech performs a synthesis of the requested text and returns the audio stream containing the speech.
func (client *Polly) SynthesizeSpeech(options SpeechOptions) (*SpeechResponse, error) {
	b, err := json.Marshal(options)
	s := string(b)
	fmt.Println(s)
	if err != nil {
		return nil, err
	}

	r, _ := http.NewRequest("POST", synthesizeSpeechAPI, bytes.NewReader(b))
	r.Header.Set("Content-Type", "application/json")

	awsClient := aws4.Client{Keys: &aws4.Keys{
		AccessKey: client.AccessKey,
		SecretKey: client.SecretKey,
	}}

	resp, err := awsClient.Do(r)
	if err != nil {
		fmt.Println("bad response")
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Could not read resp.Body")
		return nil, err

	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got non 200 status code: %s %q", resp.Status, data)
	}
	//	fmt.Println("Returning Speech Response")
	return &SpeechResponse{
		Audio:       data,
		RequestID:   resp.Header["X-Amzn-Requestcharacters"][0], // this is incorrect in the AWS polly documentation
		ContentType: resp.Header["Content-Type"][0],
	}, nil
}

// DescribeVoices retrieves list of voices from the api
// TODO handle nextToken iteration in a while loop
//http://docs.aws.amazon.com/polly/latest/dg/API_DescribeVoices.html
func (client *Polly) DescribeVoices() (*VoicesResponse, error) {

	r, err0 := http.NewRequest("GET", describeVoicesAPI, nil)
	if err0 != nil {
		fmt.Println("Get Request Error")
		return nil, err0
	}

	r.Header.Set("Content-Type", "application/json")
	awsClient := aws4.Client{Keys: &aws4.Keys{
		AccessKey: client.AccessKey,
		SecretKey: client.SecretKey,
	}}

	resp, err := awsClient.Do(r)
	fmt.Println(resp)
	fmt.Println("Got response")
	if err != nil {
		fmt.Println("Response error")
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got non 200 status code: %s %q", resp.Status, data)
	}

	voices := new(VoicesResponse)
	err = json.Unmarshal(data, voices)
	if err != nil {
		return nil, err
	}

	return voices, nil
}
