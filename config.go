package main

import (
	"github.com/spf13/viper"
	"math"
	"strconv"
	"strings"
)

type Config struct {
	N          float64
	lambda     float64
	pv         float64
	alpha      float64
	beta       float64
	k          float64
	yita       []float64
	nn         []float64
	P_v        []float64
	nodeNumber int
	maxPosX    float64
	maxPosY    float64
	receiveR   float64 // instead of beta
	timeSlot   int     // 单位 毫秒
}

func ConfigInitial() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	Conf.k = viper.GetFloat64("k")
	Conf.N = viper.GetFloat64("N")
	Conf.lambda = viper.GetFloat64("lambda")
	Conf.alpha = viper.GetFloat64("alpha")
	Conf.nodeNumber = viper.GetInt("nodeNumber")
	Conf.pv = viper.GetFloat64("pv")
	Conf.beta = viper.GetFloat64("beta")
	Conf.maxPosX = viper.GetFloat64("maxPosX")
	Conf.maxPosY = viper.GetFloat64("maxPosY")
	Conf.receiveR = viper.GetFloat64("receiveR")
	Conf.timeSlot = viper.GetInt("timeSlot")
	yitami := viper.GetInt("yitami")
	nnmi := viper.GetInt("nnmi")
	yitastr := viper.GetString("yita")
	nnstr := viper.GetString("nn")
	yitaarr := strings.Split(yitastr, ",")
	for _, v := range yitaarr {
		tmp, _ := strconv.ParseFloat(v, 64)
		Conf.yita = append(Conf.yita, tmp*math.Pow10(yitami))
	}
	nnarr := strings.Split(nnstr, ",")
	for _, v := range nnarr {
		tmp, _ := strconv.ParseFloat(v, 64)
		Conf.nn = append(Conf.nn, tmp*math.Pow10(nnmi))
	}
	Conf.P_v = append(Conf.P_v, math.Pow(100*math.Sqrt2, Conf.alpha)*Conf.beta*Conf.N)
	Conf.P_v = append(Conf.P_v, math.Pow(200, Conf.alpha)*Conf.beta*Conf.N)

	return nil
}
