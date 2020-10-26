package pingmesh_agent

import (
	"bytes"
	"fmt"
	"k8s.io/klog/v2"
	"os/exec"
	"strconv"
	"strings"
	"time"
)


type ProberResultOne struct {
	WorkerName           string
	MetricName           string
	TargetAddr           string
	SourceRegion         string
	TargetRegion         string
	ProbeType            string
	TimeStamp            int64
	Value                float32
}

func execCmd(cmdStr string) (success bool, outStr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		klog.Error("execCmdMsg ", err, " cmd ", cmdStr)

		if strings.Contains(err.Error(), "killed") {
			return false, "killed"
		}

		return false, string(stderr.Bytes())
	}
	outStr = string(stdout.Bytes())
	fmt.Println(outStr)
	return true, outStr

}

func ProbeICMPs(pl pinglist) ([]([]*ProberResultOne)){
	ts := make([]*Target,0)
	for patition,urls := range pl.PingList{
		for _, url := range urls {
			t := &Target{
				WorkName: pl.WorkerName,
				TargetAddr: url,
				SourceRegion: pl.Patition,
				TargetRegion: patition,
				ProbeType: "Ping",
			}
			ts = append(ts,t)
		}
	}

	pros := make([]([]*ProberResultOne),0)
	for _, t := range ts{
		pro := ProbeICMP(t)
		pros = append(pros, pro)
	}
	return pros
}
func ProbeICMP(t *Target) ([]*ProberResultOne){

	defer func() {
		if r := recover(); r != nil {
			//resultErr, _ := r.(error)
			klog.Errorf("ProbeICMP panic.....")

		}
	}()

	pingCmd := fmt.Sprintf("/usr/bin/timeout --signal=KILL 15s  /usr/bin/ping -q -A -f -s 100 -W 1000 -c 50 %s", t.TargetAddr)
	//level.Info(lt.logger).Log("msg", "LocalTarget  ProbeICMP start ...", "uid", lt.Uid(), "pingcmd", pingCmd)
	klog.V(6).Infof("ProbeICMP start, targetUrl:"+t.TargetAddr)
	success, outPutStr := execCmd(pingCmd)
	prs := make([]*ProberResultOne, 0)
	var (
		pkgdLine    string
		latenLinke  string
		pkgRateNum  float64
		pingEwmaNum float64
		pingSuccess float64
	)

	pkgRateNum = -1
	pingEwmaNum = -1
	pingSuccess = 0
	prSu := ProberResultOne{
		MetricName:   MetricsNamePingTargetSuccess,
		WorkerName:   t.WorkName,
		TargetAddr:   t.TargetAddr,
		SourceRegion: t.SourceRegion,
		TargetRegion: t.TargetRegion,
		ProbeType:    t.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pingSuccess),
	}
	if success == false {
		klog.V(6).Info("ProbeICMP failed, err_str: ", outPutStr)

		if strings.Contains(outPutStr, "killed") {
			klog.V(6).Info("ProbeICMP killed, err_str: ", outPutStr)
			prSu.Value = -1
			prs = append(prs, &prSu)
			return prs

		}
		return prs
	}

	for _, line := range (strings.Split(outPutStr, "\n")) {
		if strings.Contains(line, "packets transmitted") {
			pkgdLine = line
			continue
		}
		if strings.Contains(line, "min/avg/max/mdev") {
			latenLinke = line
			continue
		}
	}

	if len(pkgdLine) > 0 {

		pkgRate := strings.Split(pkgdLine, " ")[5]
		pkgRate = strings.Replace(pkgRate, "%", "", -1)
		pkgRateNum, _ = strconv.ParseFloat(pkgRate, 64)
	}

	if len(latenLinke) > 0 {
		pingEwmas := strings.Split(latenLinke, " ")

		pingEwma := pingEwmas[len(pingEwmas)-2]
		pingEwma = strings.Split(pingEwma, "/")[1]
		pingEwmaNum, _ = strconv.ParseFloat(pingEwma, 64)
	}

	klog.V(6).Infof( "ProbeICMP_one_res, pingcmd:%s, outPutStr:%s, pkgRateNum:%f, pingEwmaNum:%f",  pingCmd, outPutStr, float32(pkgRateNum), float32(pingEwmaNum))

	prDr := ProberResultOne{
		MetricName:   MetricsNamePingPackageDrop,
		WorkerName:   t.WorkName,
		TargetAddr:   t.TargetAddr,
		SourceRegion: t.SourceRegion,
		TargetRegion: t.TargetRegion,
		ProbeType:    t.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pkgRateNum),
	}

	prLaten := ProberResultOne{
		MetricName:   MetricsNamePingLatency,
		WorkerName:   t.WorkName,
		TargetAddr:   t.TargetAddr,
		SourceRegion: t.SourceRegion,
		TargetRegion: t.TargetRegion,
		ProbeType:    t.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pingEwmaNum),
	}

	if pkgRateNum == 100 {
		prSu.Value = -1
	} else {
		prSu.Value = 1
	}

	prs = append(prs, &prSu)
	prs = append(prs, &prDr)
	prs = append(prs, &prLaten)

	return prs
}