package user_agent

import (
	"os"
	"log"
	"bufio"
	"path/filepath"

	uaparser "../uaparser"
	lru "github.com/hashicorp/golang-lru"
	"sync"
)

type Configuration struct {
	cacheSize int
	dataDir string
}

var config = Configuration{
	cacheSize: 10000,
	dataDir: "./data",
}

type UserAgent struct {
	parser *uaparser.Parser
	lruCache *lru.Cache
}

func New() (*UserAgent, error) {
	//data, err := ioutil.ReadFile(file)
	//if err != nil {
	//	log.Fatal(err)
	//	return nil, err
	//}
	//
	//var config *Configuration = &Configuration{}
	//if err := yaml.Unmarshal(data, config); err != nil {
	//	log.Fatal(err)
	//	return nil, err
	//}

	//main folder
	pwd, err := os.Getwd()

	// init parser instance
	parser, err := uaparser.New(filepath.Join(pwd, "..", "uap-core",  "regexes.yaml"))
	if err != nil {
		log.Fatal("Error read regexp file:", err)
		return nil, err
	}

	// init lru cache
	lruCache, err := lru.New(config.cacheSize)
	if err != nil {
		log.Fatal("Error init lru:", err)
		return nil, err
	}

	ua := &UserAgent{
		parser: parser,
		lruCache: lruCache,
	}

	// warm up lru cache
	ua.warmupCache()

	return ua, nil
}

func (ua *UserAgent) warmupCache() {
	log.Println("Start warm up cache...")
	filename, _ := filepath.Abs(config.dataDir)

	datafile := filepath.Join(filename, "ua_warmup.log")
	file, err := os.Open(datafile)
	if err != nil {
		log.Println("Cache was not warming up:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ua.parseUAWithLRU(scanner.Text())
	}
}

func (ua *UserAgent) parseUAWg(user_agent string) *uaparser.Client {
	var wg sync.WaitGroup

	wg.Add(3)
	cli := &uaparser.Client{}

	go func() {
		cli.UserAgent = ua.parser.ParseUserAgent(user_agent)
		wg.Done()
	}()

	go func() {
		cli.Os = ua.parser.ParseOs(user_agent)
		wg.Done()
	}()

	go func() {
		cli.Device = ua.parser.ParseDevice(user_agent)
		wg.Done()
	}()

	wg.Wait()
	return cli
}

func (ua *UserAgent) parseUA(user_agent string) *uaparser.Client {
	return ua.parser.Parse(user_agent)
}

func (ua *UserAgent) parseUAWithLRU(user_agent string) *uaparser.Client {
	if val, ok := ua.lruCache.Get(user_agent); ok {
		return val.(*uaparser.Client)
	}

	c := ua.parseUAWg(user_agent)
	ua.lruCache.Add(user_agent, c)
	return c
}

func (ua *UserAgent) parseUAWithLRUWithoutWG(user_agent string) *uaparser.Client {
	if val, ok := ua.lruCache.Get(user_agent); ok {
		return val.(*uaparser.Client)
	}

	c := ua.parseUA(user_agent)
	ua.lruCache.Add(user_agent, c)
	return c
}