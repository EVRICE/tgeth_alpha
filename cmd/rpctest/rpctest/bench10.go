package rpctest

import (
	"fmt"
	"net/http"
	"time"
)

func Bench10(tgUrl, gethUrl string, bn uint64) error {
	var client = &http.Client{
		Timeout: time.Second * 600,
	}

	var res CallResult
	reqGen := &RequestGenerator{
		client: client,
	}

	reqGen.reqID++

	var b EthBlockByNumber
	res = reqGen.TurboGeth("eth_getBlockByNumber", reqGen.getBlockByNumber(int(bn)), &b)
	if res.Err != nil {
		return fmt.Errorf("retrieve block (turbo-geth) %d: %v", bn, res.Err)
	}
	if b.Error != nil {
		return fmt.Errorf("retrieving block (turbo-geth): %d %s", b.Error.Code, b.Error.Message)
	}
	for _, tx := range b.Result.Transactions {
		reqGen.reqID++

		var trace EthTxTrace
		res = reqGen.TurboGeth("debug_traceTransaction", reqGen.traceTransaction(tx.Hash), &trace)
		if res.Err != nil {
			fmt.Printf("Could not trace transaction (turbo-geth) %s: %v\n", tx.Hash, res.Err)
			print(client, routes[TurboGeth], reqGen.traceTransaction(tx.Hash))
		}

		if trace.Error != nil {
			fmt.Printf("Error tracing transaction (turbo-geth): %d %s\n", trace.Error.Code, trace.Error.Message)
		}

		var traceg EthTxTrace
		res = reqGen.Geth("debug_traceTransaction", reqGen.traceTransaction(tx.Hash), &traceg)
		if res.Err != nil {
			print(client, routes[Geth], reqGen.traceTransaction(tx.Hash))
			return fmt.Errorf("trace transaction (geth) %s: %v", tx.Hash, res.Err)
		}
		if traceg.Error != nil {
			return fmt.Errorf("tracing transaction (geth): %d %s", traceg.Error.Code, traceg.Error.Message)
		}
		if res.Err == nil && trace.Error == nil {
			if !compareTraces(&trace, &traceg) {
				return fmt.Errorf("different traces block %d, tx %s", bn, tx.Hash)
			}
		}
		reqGen.reqID++
	}
	return nil
}
