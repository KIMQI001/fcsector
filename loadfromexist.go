package fcsector

import (
	"context"
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	commcid "github.com/filecoin-project/go-fil-commcid"
	ffi "github.com/filecoin-project/lotus/extern/filecoin-ffi"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

type sizeToInfo struct {
	size int64
	path string
	commP []byte
}
var (
	sizeOf512MiB = sizeToInfo{
		//size:512 << 20,
		size:532676608,
		path:"",
		commP:nil,
	}
	SizeOf1GiB = sizeToInfo{
		size:1 << 30,
		path:"",
		commP:nil,
	}
	SizeOf32GiB = sizeToInfo{
		//size:32 << 30,
		size:34091302912,
		path:"",
		commP:nil,
	}
)

type existFileInfo struct {
	Path string `json:"path"`
	CommP []byte `json:"comm_p"`
}
func (sb *SectorBuilder) AddPiece(ctx context.Context, pieceSize abi.UnpaddedPieceSize, sectorNum abi.SectorNumber, file io.Reader, existingPieceSizes []abi.UnpaddedPieceSize) (abi.PieceInfo, error) {
	// add at the head
	if existCID :=checkExistFile(pieceSize);len(existCID)>30{
		log.Infof("\n-----the size of piece is exist!!-----\npiece size is %d\ncommp is %d \nthe cid is %s\n",pieceSize,existCID,commcid.PieceCommitmentV1ToCID(existCID[:]))
		return abi.PieceInfo{
			Size:     pieceSize.Padded(),
			PieceCID: commcid.PieceCommitmentV1ToCID(existCID[:]),
		},nil
	}
	// add at the second last
	comp,_:=commcid.CIDToDataCommitmentV1(pieceCID)
	log.Infof("the piece CID is %s\n the commp is %s",pieceCID,strings.Replace(strings.Trim(fmt.Sprint(comp), "[]"), " ", ",", -1))

}
func checkExistFile(pieceSize abi.UnpaddedPieceSize) []byte{
	// check the pieceSize
	log.Info("start check!",pieceSize)
	log.Infof("512 is %d,1GB is %d,32Gb is %d",sizeOf512MiB.size,SizeOf1GiB.size,SizeOf32GiB.size)
	if int64(pieceSize)==sizeOf512MiB.size{
		return sizeOf512MiB.commP
	}else if int64(pieceSize)==SizeOf1GiB.size{
		return SizeOf1GiB.commP
	}else if int64(pieceSize)==SizeOf32GiB.size{
		return SizeOf32GiB.commP
	}
	// if not exist,continue
	return nil
}
// TODO:this should embed in lotus-miner-run
func ExistInfoInit()  {
	var existInfos map[string]existFileInfo
	if err := json.Unmarshal(existFileInfoFetch(),&existInfos);err!=nil{
		return
	}
	for size, info := range existInfos{
		if strings.Contains(size,"512"){
			sizeOf512MiB.path = info.Path
			sizeOf512MiB.commP = info.CommP
		}else if strings.Contains(size,"1GiB"){
			SizeOf1GiB.path = info.Path
			SizeOf1GiB.commP = info.CommP
		}else if strings.Contains(size,"32GiB"){
			SizeOf32GiB.path = info.Path
			SizeOf32GiB.commP = info.CommP
		}
	}
}

func existFileInfoFetch() []byte {
	return rice.MustFindBox("existinfo").MustBytes("existFileInfo.json")
}
