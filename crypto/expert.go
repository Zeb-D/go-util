package crypto

import (
	"bytes"
	"encoding/json"
	"errors"
)

//云端拉取策略请求
type GetDataFromAtopRequest struct {
	SecKey     string
	TaskId     string
	GwId       string
	TransId    string
	BizType    int
	DataType   int
	PageNo     int
	PageSize   int
	Result     int
	FailedData []int
	ErrMsg     string
}

type StrategiesOrJobSendResp struct {
	Result    struct{ Result string } `json:"result"`
	T         int64                   `json:"t"`
	Success   bool                    `json:"success"`
	ErrorCode string                  `json:"errorCode"`
}

func strategiesOrJobResp(atopUrl, gwId, secKey string, req GetDataFromAtopRequest) (StrategiesOrJobSendResp, error) {
	var (
		err      error
		respBody []byte
		resp     StrategiesOrJobSendResp
	)

	params := map[string]interface{}{
		"task_id":     req.TaskId,
		"gw_id":       req.GwId,
		"trans_id":    req.TransId,
		"biz_type":    req.BizType,
		"data_type":   req.DataType,
		"result":      req.Result,
		"failed_data": req.FailedData,
		"err_msg":     req.ErrMsg,
	}

	if respBody, err = GatewayRequestWithMap(
		atopUrl,
		"xxx.xxx.result.report",
		gwId,
		"1.0",
		secKey,
		params); err != nil {
		return StrategiesOrJobSendResp{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(respBody))
	decoder.UseNumber()
	if err = decoder.Decode(&resp); err != nil {
		return StrategiesOrJobSendResp{}, err
	}
	if !resp.Success {
		return resp, errors.New(resp.ErrorCode)
	}
	return resp, nil
}

//
//func getStrategy(atopUrl, gwId, secKey string, req GetDataFromAtopRequest) (dtos.StrategiesFromAtopResp, error) {
//	var (
//		err      error
//		respBody []byte
//		resp     dtos.StrategiesFromAtopResp
//	)
//
//	params := map[string]interface{}{
//		"task_id":   req.TaskId,
//		"gw_id":     req.GwId,
//		"trans_id":  req.TransId,
//		"biz_type":  req.BizType,
//		"data_type": req.DataType,
//		"page_no":   req.PageSize,
//		"page_size": req.PageSize,
//	}
//
//	if respBody, err = GatewayRequestWithMap(
//		atopUrl,
//		AtopStrategyOrJobGet,
//		gwId,
//		"1.0",
//		secKey,
//		params); err != nil {
//		return dtos.StrategiesFromAtopResp{}, err
//	}
//
//	decoder := json.NewDecoder(bytes.NewReader(respBody))
//	decoder.UseNumber()
//	if err = decoder.Decode(&resp); err != nil {
//		return dtos.StrategiesFromAtopResp{}, err
//	}
//	if !resp.Success {
//		return resp, errors.New(resp.ErrorCode)
//	}
//	return resp, nil
//}
//
//func getJob(atopUrl, gwId, secKey string, req GetDataFromAtopRequest) (dtos.JobsFromAtopResp, error) {
//	var (
//		err      error
//		respBody []byte
//		resp     dtos.JobsFromAtopResp
//	)
//
//	params := map[string]interface{}{
//		"task_id":   req.TaskId,
//		"gw_id":     req.GwId,
//		"trans_id":  req.TransId,
//		"biz_type":  req.BizType,
//		"data_type": req.DataType,
//		"page_no":   req.PageSize,
//		"page_size": req.PageSize,
//	}
//
//	if respBody, err = GatewayRequestWithMap(
//		atopUrl,
//		AtopStrategyOrJobGet,
//		gwId,
//		"1.0",
//		secKey,
//		params); err != nil {
//		return dtos.JobsFromAtopResp{}, err
//	}
//
//	decoder := json.NewDecoder(bytes.NewReader(respBody))
//	decoder.UseNumber()
//	if err = decoder.Decode(&resp); err != nil {
//		return dtos.JobsFromAtopResp{}, err
//	}
//	if !resp.Success {
//		return resp, errors.New(resp.ErrorCode)
//	}
//	return resp, nil
//}
//
//func getStrategyOrJobIds(atopUrl, gwId, secKey string, req GetDataFromAtopRequest) (dtos.StrategyOrJobIdsFromAtopResp, error) {
//	var (
//		err      error
//		respBody []byte
//		resp     dtos.StrategyOrJobIdsFromAtopResp
//	)
//
//	params := map[string]interface{}{
//		"task_id":   req.TaskId,
//		"gw_id":     req.GwId,
//		"trans_id":  req.TransId,
//		"biz_type":  req.BizType,
//		"data_type": req.DataType,
//		"page_no":   req.PageSize,
//		"page_size": req.PageSize,
//	}
//
//	if respBody, err = GatewayRequestWithMap(
//		atopUrl,
//		AtopStrategyOrJobGet,
//		gwId,
//		"1.0",
//		secKey,
//		params); err != nil {
//		return dtos.StrategyOrJobIdsFromAtopResp{}, err
//	}
//
//	decoder := json.NewDecoder(bytes.NewReader(respBody))
//	decoder.UseNumber()
//	if err = decoder.Decode(&resp); err != nil {
//		return StrategyOrJobIdsFromAtopResp{}, err
//	}
//	if !resp.Success {
//		return resp, errors.New(resp.ErrorCode)
//	}
//	return resp, nil
//}
