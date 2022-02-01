package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	c1 "github.com/sim15/anon-com/construction1"
	"github.com/sim15/anon-com/sposs"
)

var group, q = c1.DefaultSetup()

var pp = sposs.NewPublicParams(group)

func main() {

	// NUM_ACCOUNTS := []uint64{1 << 2, 1 << 5, 1 << 8, 1 << 11, 1 << 14, 1 << 17, 1 << 18, 1 << 19, 1 << 20}
	NUM_ACCOUNTS := []uint64{1 << 14}
	// MESSAGE_SIZE := []uint64{50, 100, 300, 500, 750, 1000}
	MESSAGE_SIZE := []uint64{1000}
	NUM_TRIALS := 10

	// var experiment Experiment
	exps := make([]Experiment, 0)
	for _, n := range NUM_ACCOUNTS {

		for _, l := range MESSAGE_SIZE {

			experiment := &Experiment{
				NumBoxes:        n,
				MLength:         l,
				Construction1MS: make([][]int64, 3),
			}

			// construction1MS = append(construction1MS, make([][]int64, 1))

			for trial := 0; trial < NUM_TRIALS; trial++ {
				trialTime := benchmarkConst1(n, l)
				experiment.Construction1MS[0] = append(experiment.Construction1MS[0], trialTime[0])
				experiment.Construction1MS[1] = append(experiment.Construction1MS[1], trialTime[1])
				experiment.Construction1MS[2] = append(experiment.Construction1MS[2], trialTime[2])
				// experiment.Construction1MS = append(experiment.Construction1MS)
				fmt.Printf("Finished trial %v/%v\n", trial, NUM_TRIALS)
			}
			exps = append(exps, *experiment)

		}
		fmt.Printf("Finished database size of %v mailboxes\n", n)
	}
	fmt.Printf("Finished %v trials of %v database sizes and %v message sizes\n", NUM_TRIALS, len(NUM_ACCOUNTS), len(MESSAGE_SIZE))

	experimentJSON, _ := json.MarshalIndent(exps, "", " ")
	ioutil.WriteFile("experiment.json", experimentJSON, 0644)
}

func benchmarkConst1(numMailboxes, messageSize uint64) []int64 {
	timing := make([]int64, 3)

	serverTestX := pp.ExpField.RandomElement().Int
	serverTestAltX := pp.ExpField.Add(pp.ExpField.NewElement(serverTestX), q).Int

	sA := c1.NewSpossServer(false, messageSize, numMailboxes, pp)
	sB := c1.NewSpossServer(true, messageSize, numMailboxes, pp)

	sA.InitTestSpossSlotList(group.NewElement(serverTestX))
	sB.InitTestSpossSlotList(group.NewElement(serverTestAltX))

	x := pp.ExpField.RandomElement().Int
	altX := pp.ExpField.Add(pp.ExpField.NewElement(x), q).Int

	sA.Boxes.Slots[0].SPOSSKey = pp.Group.NewElement(x)
	sB.Boxes.Slots[0].SPOSSKey = pp.Group.NewElement(altX)

	sA.SPOSSKeys[0] = pp.Group.NewElement(x)
	sB.SPOSSKeys[0] = pp.Group.NewElement(altX)

	message := make([]byte, messageSize)
	rand.Read(message)
	c := c1.NewClient(pp, x, messageSize, message)
	// -----------------------------------
	start := time.Now()

	query := c.NewClientQuery(0, numMailboxes, x) // TODO: nummailboxes will be constant in practice, so fix later

	timing[0] = time.Since(start).Milliseconds()
	// -----------------------------------

	sB.StartSession(query[1])

	sB.ComputePrepareAuthAudit()
	rand := sA.Boxes.ProofParams.Group.Field.RandomElement()
	sB.Boxes.ProofParams.SetRandSeed(rand)

	pubAuditShareB, _ := sB.Boxes.ProofParams.PrepareAudit(
		sB.CurrentSession.RecievedQuery.SPoSSProof,
		sB.CurrentSession.QueryShare.Share,
		true)

	// -----------------------------
	start = time.Now()

	sA.StartSession(query[0])
	sA.ComputePrepareAuthAudit()

	sA.Boxes.ProofParams.SetRandSeed(rand)

	_, privAuditShareA := sA.Boxes.ProofParams.PrepareAudit(
		sA.CurrentSession.RecievedQuery.SPoSSProof,
		sA.CurrentSession.QueryShare.Share,
		false)

	pubVerificationShareA, privVerificationShareA := sA.Boxes.ProofParams.Audit(
		pubAuditShareB, //for benchmark
		privAuditShareA,
		false)

	sA.Boxes.ProofParams.VerifyAudit(pubVerificationShareA, privVerificationShareA)

	timing[1] = time.Since(start).Milliseconds()

	// -----------------------------------

	start = time.Now()

	sA.WriteShare()

	timing[2] = time.Since(start).Milliseconds()

	return timing

}

func benchmarkExpress(numMailboxes, messageSize uint64) {
	//  TODO: only consider messages of size 1000 bytes (1kb)
}
