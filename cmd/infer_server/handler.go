package main

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/inference"
	"github.com/ethereum/go-ethereum/inference/synapse"
	"github.com/ethereum/go-ethereum/log"
)

func infoHashHandler(w http.ResponseWriter, inferWork *inference.IHWork) {
	if inferWork.Model == "" {
		RespErrorText(w, ErrModelEmpty)
		return
	}
	if inferWork.Input == "" {
		RespErrorText(w, ErrInputEmpty)
		return
	}

	log.Info("Infer Task", "Model Hash", inferWork.Model, "Input Hash", inferWork.Input)
	label, err := synapse.Engine().InferByInfoHash(inferWork.Model, inferWork.Input)

	if err == nil {
		log.Info("Infer Succeed", "result", label)
		RespInfoText(w, label)
	} else {
		log.Warn("Infer Failed", "error", err)
		RespErrorText(w, err)
	}
}

func inputContentHandler(w http.ResponseWriter, inferWork *inference.ICWork) {
	if inferWork.Model == "" {
		RespErrorText(w, ErrModelEmpty)
		return
	}

	model, input := inferWork.Model, inferWork.Input

	log.Info("Infer Work", "Model Hash", model, "Input Content", input)
	var cacheKey = synapse.RLPHashString(fmt.Sprintf("%s:%x", model, input))
	if v, ok := simpleCache.Load(cacheKey); ok && !(*IsNotCache) {
		log.Info("Infer succeed via cache", "cache key", cacheKey, "label", v.([]byte))
		RespInfoText(w, v.([]byte))
		return
	}

	// Fixed bugs, ctx_getSolidityBytes returns 0x which stands for state invalid
	if len(input) == 0 {
		log.Warn("Input content state invalid", "error", "bytes length is zero")
		RespErrorText(w, "input bytes length is zero")
		return
	}

	label, err := synapse.Engine().InferByInputContent(model, input)

	if err != nil {
		log.Warn("Infer Failed", "error", err)
		RespErrorText(w, err)
		return
	}

	log.Info("Infer Succeed", "result", label)
	if !(*IsNotCache) {
		simpleCache.Store(cacheKey, label)
	}

	RespInfoText(w, label)
}