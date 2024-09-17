package main

import (
	"os"
	"errors"
	"strconv"
	"fmt"

	"github.com/urfave/cli/v2"
	deserializer "github.com/worldcoin/ptau-deserializer/deserialize"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
	"github.com/worldcoin/semaphore-mtb-setup/keys"
	"github.com/worldcoin/semaphore-mtb-setup/phase1"
	"github.com/worldcoin/semaphore-mtb-setup/phase2"
)


func p1n(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	powerStr := cCtx.Args().Get(0)
	power, err := strconv.Atoi(powerStr)
	if err != nil {
		return err
	}
	if power > 26 {
		return errors.New("can't support powers larger than 26")
	}
	outputPath := cCtx.Args().Get(1)
	err = phase1.Initialize(byte(power), outputPath)
	return err
}

func p1c(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	outputPath := cCtx.Args().Get(1)
	err := phase1.Contribute(inputPath, outputPath)
	return err
}

func p1i(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}

	ptauFilePath := cCtx.Args().Get(0)
	outputFilePath := cCtx.Args().Get(1)

	ptau, err := deserializer.ReadPtau(ptauFilePath)
	if err != nil {
		return err
	}

	phase1, err := deserializer.ConvertPtauToPhase1(ptau)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	_, err = phase1.WriteTo(outputFile)
	if err != nil {
		return err
	}

	return nil
}

// func p2n(cCtx *cli.Context) error {
// 	// sanity check
// 	if cCtx.Args().Len() != 3 {
// 		return errors.New("please provide the correct arguments")
// 	}

// 	phase1Path := cCtx.Args().Get(0)
// 	r1csPath := cCtx.Args().Get(1)
// 	phase2Path := cCtx.Args().Get(2)
// 	err := phase2.Initialize(phase1Path, r1csPath, phase2Path)
// 	return err
// }

func p2n(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 4 {
		return errors.New("please provide the correct arguments")
	}

	phase1Path := cCtx.Args().Get(0)
	r1csPath := cCtx.Args().Get(1)
	phase2Path := cCtx.Args().Get(2)
	evalsPath := cCtx.Args().Get(3)

	phase1File, err := os.Open(phase1Path)
	if err != nil {
		return err
	}
	phase1 := &mpcsetup.Phase1{}
	phase1.ReadFrom(phase1File)

	r1csFile, err := os.Open(r1csPath)
	if err != nil {
		return err
	}
	r1cs := cs.R1CS{}
	r1cs.ReadFrom(r1csFile)

	phase2, evals := mpcsetup.InitPhase2(&r1cs, phase1)

	phase2File, err := os.Create(phase2Path)
	if err != nil {
		return err
	}
	phase2.WriteTo(phase2File)

	evalsFile, err := os.Create(evalsPath)
	if err != nil {
		return err
	}
	evals.WriteTo(evalsFile)

	return nil
}

// func Initialize(phase1Path, r1csPath, phase2Path string) error {
// 	phase1File, err := os.Open(phase1Path)
// 	if err != nil {
// 		return err
// 	}
// 	defer phase1File.Close()

// 	phase2File, err := os.Create(phase2Path)
// 	if err != nil {
// 		return err
// 	}
// 	defer phase2File.Close()

// 	// 1. Process Headers
// 	header1, header2, err := processHeader(r1csPath, phase1File, phase2File)
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Convert phase 1 SRS to Lagrange basis
// 	if err := processLagrange(header1, header2, phase1File, phase2File); err != nil {
// 		return err
// 	}

// 	// 3. Process evaluation
// 	if err := processEvaluations(header1, header2, r1csPath, phase1File); err != nil {
// 		return err
// 	}

// 	// Evaluate Delta and Z
// 	if err := processDeltaAndZ(header1, header2, phase1File, phase2File); err != nil {
// 		return err
// 	}

// 	// Process parameters
// 	if err := processPVCKK(header1, header2, r1csPath, phase2File); err != nil {
// 		return err
// 	}

// 	fmt.Println("Phase 2 has been initialized successfully")
// 	return nil
// }

// func p1t(cCtx *cli.Context) error {
// 	// sanity check
// 	if cCtx.Args().Len() != 4 {
// 		return errors.New("please provide the correct arguments")
// 	}
// 	inputPath := cCtx.Args().Get(0)
// 	outputPath := cCtx.Args().Get(1)
// 	inPowStr := cCtx.Args().Get(2)
// 	inPower, err := strconv.Atoi(inPowStr)
// 	if err != nil {
// 		return err
// 	}
// 	outPowStr := cCtx.Args().Get(3)
// 	outPower, err := strconv.Atoi(outPowStr)
// 	if err != nil {
// 		return err
// 	}
// 	if inPower < outPower {
// 		return errors.New("cannot transform to a higher power")
// 	}
// 	err = phase1.Transform(inputPath, outputPath, byte(inPower), byte(outPower))
// 	return err
// }

func p1v(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 1 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	err := phase1.Verify(inputPath, "")
	return err
}

func p1vt(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	transformedPath := cCtx.Args().Get(1)
	err := phase1.Verify(inputPath, transformedPath)
	return err
}

func p2c(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	outputPath := cCtx.Args().Get(1)
	err := phase2.Contribute(inputPath, outputPath)
	return err
}

func p2v(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	originPath := cCtx.Args().Get(1)
	err := phase2.Verify(inputPath, originPath)
	return err
}

func extract(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 1 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	err := keys.ExtractKeys(inputPath)
	return err
}

func exportSol(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 1 {
		return errors.New("please provide the correct arguments")
	}
	session := cCtx.Args().Get(0)
	err := keys.ExportSol(session)
	return err
}
