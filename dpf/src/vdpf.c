#include "../include/dpf.h"
#include "../include/mmo.h"
#include "../include/common.h"
#include <openssl/rand.h>

void genVDPF(
	EVP_CIPHER_CTX *ctx, 
	struct Hash *hash, 
	int size, 
	uint64_t index, 
	unsigned char* k0, 
	unsigned char *k1) {

	int didFinish = false; 
	while (!didFinish) { 

		uint128_t seeds0[size+1];
		uint128_t seeds1[size+1];
		int bits0[size + 1];
		int bits1[size + 1];

		uint128_t sCW[size];
		int tCW0[size];
		int tCW1[size];

		seeds0[0] = getRandomBlock();
		seeds1[0] = getRandomBlock();
		bits0[0] = 0;
		bits1[0] = 1;

		uint128_t s0[2], s1[2]; // 0=L,1=R
		int t0[2], t1[2];
		for(int i = 1; i <= size; i++){
			dpfPRG(ctx, seeds0[i-1], &s0[LEFT], &s0[RIGHT], &t0[LEFT], &t0[RIGHT]);
			dpfPRG(ctx, seeds1[i-1], &s1[LEFT], &s1[RIGHT], &t1[LEFT], &t1[RIGHT]);

			int keep, lose;
			int indexBit = getbit(index, size, i);
			if (indexBit == 0){
				keep = LEFT;
				lose = RIGHT;
			} else{
				keep = RIGHT;
				lose = LEFT;
			}

			sCW[i-1] = s0[lose] ^ s1[lose];

			tCW0[i-1] = t0[LEFT] ^ t1[LEFT] ^ indexBit ^ 1;
			tCW1[i-1] = t0[RIGHT] ^ t1[RIGHT] ^ indexBit;

			if (bits0[i-1] == 1){
				seeds0[i] = s0[keep] ^ sCW[i-1];
				if (keep == 0)
					bits0[i] = t0[keep] ^ tCW0[i-1];
				else 
					bits0[i] = t0[keep] ^ tCW1[i-1];
			} else{
				seeds0[i] = s0[keep];
				bits0[i] = t0[keep];
			}

			if (bits1[i-1] == 1){
				seeds1[i] = s1[keep] ^ sCW[i-1];
				if (keep == 0)
					bits1[i] = t1[keep] ^ tCW0[i-1];
				else
					bits1[i] = t1[keep] ^ tCW1[i-1];
			} else{
				seeds1[i] = s1[keep];
				bits1[i] = t1[keep];
			}
		}

		// *********************************
		// START: DPF verification code 
		// *********************************
       
	    uint128_t hashinput[2];
        hashinput[0] = index;
        hashinput[1] = seeds0[size];

		uint128_t pi0[hash->outblocks];
		uint128_t pi1[hash->outblocks];


		// printf(" in(%i) = %llx||%llx\n", 0, (unsigned long long)hashinput[0], (unsigned long long)hashinput[1]);
        mmoHash2to4(hash, (uint8_t*)&hashinput[0], (uint8_t*)&pi0);
	    // printf(" out(%i) = %llx||%llx||%llx||%llx\n", 0, (unsigned long long)pi0[0], (unsigned long long)pi0[1], (unsigned long long)pi0[2], (unsigned long long)pi0[3]);

		hashinput[0] = index;
        hashinput[1] = seeds1[size];
		// printf(" in(%i) = %llx||%llx\n", 1, (unsigned long long)hashinput[0], (unsigned long long)hashinput[1]);
        mmoHash2to4(hash, (uint8_t*)&hashinput[0], (uint8_t*)&pi1);
	    // printf(" out(%i) = %llx||%llx||%llx||%llx\n", 1, (unsigned long long)pi1[0], (unsigned long long)pi1[1], (unsigned long long)pi1[2], (unsigned long long)pi1[3]);

		uint128_t cs[4];
		cs[0] = pi0[0]^pi1[0];
		cs[1] = pi0[1]^pi1[1];
		cs[2] = pi0[2]^pi1[2];
		cs[3] = pi0[3]^pi1[3];

	    // printf(" cs(%i) = %llx||%llx||%llx||%llx\n", (int)hash->outblocks,  (unsigned long long)cs[0], (unsigned long long)cs[1], (unsigned long long)cs[2], (unsigned long long)cs[3]);

        int bit0 = seed_lsb(seeds0[size]);
        int bit1 = seed_lsb(seeds1[size]);

        if (bit0 != bit1)
            didFinish = true;
        else continue;

		// printf("seeds[size] = %llx\n", (unsigned long long)seeds0[size]);
		// printf("pi0 = %llx\n", (unsigned long long)pi0[0]);
		// printf("seeds1[size] = %llx\n", (unsigned long long)seeds1[size]);
		// printf("pi1 = %llx\n", (unsigned long long)pi1[0]);
		// *********************************
		// END: DPF verification code 
		// *********************************

        uint128_t sFinal0 = convert(&seeds0[size]);
		uint128_t sFinal1 = convert(&seeds1[size]);
		uint128_t lastCW = 1 ^ sFinal0 ^ sFinal1;

	    k0[0] = 0;
		memcpy(&k0[1], seeds0, 16);
		k0[17] = bits0[0];
		for(int i = 1; i <= size; i++){
			memcpy(&k0[18 * i], &sCW[i-1], 16);
			k0[CWSIZE * i + CWSIZE-2] = tCW0[i-1];
			k0[CWSIZE * i + CWSIZE-1] = tCW1[i-1];
		}
		memcpy(&k0[INDEX_LASTCW], &lastCW, 16);
		memcpy(&k0[INDEX_LASTCW + 16], cs, 16 * (hash->outblocks));

		memcpy(k1, k0, INDEX_LASTCW + 16 + 16 * (hash->outblocks));
		memcpy(&k1[1], seeds1, 16); // only value that is different from k0
		k1[0] = 1;
		k1[17] = bits1[0];
    }
}

// Follows implementation of https://eprint.iacr.org/2021/580.pdf (Figure 1)
// Hash1 = H, Hash2 = H' 
// pi is the verification output (should be equal on both servers)
void batchEvalVDPF(
	EVP_CIPHER_CTX *ctx, 
	struct Hash *hash1, 
	struct Hash *hash2, 
	int size,
	bool b, 
	unsigned char* k, 
	uint64_t* in, 
	uint64_t inl, 
	uint8_t* out, 
	uint8_t* proof) {
	
	// parse the key 
	uint128_t seeds[size+1];
	int bits[size+1];
	uint128_t sCW[size+1];
	int tCW0[size];
	int tCW1[size];
	uint128_t cs[4];
	uint128_t pi[4];

	memcpy(&seeds[0], &k[1], 16);
	bits[0] = b;

	for (int i = 1; i <= size; i++){
		memcpy(&sCW[i-1], &k[18 * i], 16);
		tCW0[i-1] = k[CWSIZE * i + CWSIZE-2];
		tCW1[i-1] = k[CWSIZE * i + CWSIZE-1];
	}

	memcpy(cs, &k[INDEX_LASTCW + 16], 16 * (hash1->outblocks));
	memcpy(pi, &k[INDEX_LASTCW + 16], 16 * (hash1->outblocks)); // pi = cs
	
	// printf("cs(%i) = %llx||%llx||%llx||%llx\n", b, (unsigned long long)pi[0], (unsigned long long)pi[2], (unsigned long long)pi[3], (unsigned long long)pi[3]);

	uint128_t hashinput[4];
	uint128_t tpi[4];
	uint128_t cpi[2];

	uint128_t sL, sR;
	int tL, tR;

	// outter loop: iterate over all evaluation points 
	for (int l = 0; l < inl; l++) { 

	 	// printf("start_seed = %llx\n", (unsigned long long)seeds[0]);
		for (int i = 1; i <= size; i++){

			dpfPRG(ctx, seeds[i-1], &sL, &sR, &tL, &tR);
			
			if (bits[i-1] == 1){
				sL = sL ^ sCW[i-1];
				sR = sR ^ sCW[i-1];
				tL = tL ^ tCW0[i-1];
				tR = tR ^ tCW1[i-1];
			}

			int xbit = getbit(in[l], size, i);
			
			seeds[i] = (1-xbit) * sL + xbit * sR;
			bits[i] = (1-xbit) * tL +  xbit * tR;
		}
		
		// *********************************
		// START: DPF verification code 
		// *********************************
		int bit = seed_lsb(seeds[size]);
		// printf("S%i l: %i bit %i ", b, l, bit);
		// printf(" pi = %llx", (unsigned long long)pi[0]);

        hashinput[0] = in[l];
		hashinput[1] = seeds[size];

		// step 1: H(seeds[size]||X[l])
		// printf(" in1(%i) = %llx||%llx\n", b, (unsigned long long)hashinput[0], (unsigned long long)hashinput[1]);
		mmoHash2to4(hash1, (uint8_t*)&hashinput[0], (uint8_t*)&tpi[0]);
		// printf(" out1(%i) = %llx||%llx||%llx||%llx\n", b, (unsigned long long)tpi[0], (unsigned long long)tpi[1], (unsigned long long)tpi[2], (unsigned long long)tpi[3]);

		// step 2: pi^correct(tpi, cs, bit)
		hashinput[0] = pi[0]^correct(tpi[0], cs[0], bit);
		hashinput[1] = pi[1]^correct(tpi[1], cs[1], bit);
		hashinput[2] = pi[2]^correct(tpi[2], cs[2], bit);
		hashinput[3] = pi[3]^correct(tpi[3], cs[3], bit);
	
		// step 3: H'(step2)
		// printf(" in2(%i) = %llx||%llx||%llx||%llx\n", b, (unsigned long long)hashinput[0], (unsigned long long)hashinput[1], (unsigned long long)hashinput[2], (unsigned long long)hashinput[3]);
		mmoHash4to2(hash2, (uint8_t*)&hashinput[0], (uint8_t*)&cpi[0]);
		// printf(" out2(%i) = %llx||%llx\n", b, (unsigned long long)cpi[0], (unsigned long long)cpi[1]);
		
		pi[0] ^= cpi[0]; 
		pi[1] ^= cpi[1];

		// *********************************
		// END: DPF verification code 
		// *********************************
		
		// printf("pi(%i) = %llx\n", b, (unsigned long long)pi[0]);

		uint128_t res = convert(&seeds[size]);

		if (bit == 1) {
			//correction word
			res = res ^ convert((uint128_t*)&k[18*size+18]);
		}

		// copy block to byte output
		memcpy(&out[l*sizeof(uint128_t)], &res, sizeof(uint128_t));
					
		// printf("results[l] = %llx\n", (unsigned long long)results[l]);
	}

	// printf("pi = %llx%llx%llx%llx\n", (unsigned long long)pi[0], (unsigned long long)pi[1], (unsigned long long)pi[2], (unsigned long long)pi[3]);

	memcpy(proof, pi, sizeof(uint128_t) * hash1->outblocks);
}


/* Need to allow specifying start and end for dataShare */
void fullDomainVDPF(
	EVP_CIPHER_CTX *ctx, 
	struct Hash *hash1, 
	struct Hash *hash2, 
	int size, 
	bool b, 
	unsigned char* k, 
	uint8_t *out, 
	uint8_t *proof){

    //dataShare is of size dataSize
    int numLeaves = 1 << size;
	int maxLayer = size;

    int currLevel = 0;
    int levelIndex = 0;
    int numIndexesInLevel = 2;

    int treeSize = 2 * numLeaves - 1;

 	// treesize too big to allocate on stack
	uint128_t *seeds = malloc(sizeof(uint128_t)*treeSize);
	int *bits = malloc(sizeof(int)*treeSize);
	uint128_t sCW[maxLayer+1];
	int tCW0[maxLayer+1];
	int tCW1[maxLayer+1];
	uint128_t cs[hash1->outblocks];
	uint128_t pi[hash1->outblocks];

	memcpy(seeds, &k[1], 16);
	memcpy(cs, &k[INDEX_LASTCW + 16], 16 * (hash1->outblocks));
	memcpy(pi, &k[INDEX_LASTCW + 16], 16 * (hash1->outblocks)); // pi = cs

	bits[0] = b;

	for (int i = 1; i <= maxLayer; i++){
		memcpy(&sCW[i-1], &k[18 * i], 16);
		tCW0[i-1] = k[18 * i + 16];
		tCW1[i-1] = k[18 * i + 17];
	}
	

	uint128_t sL, sR;
	int tL, tR;
    int parentIndex = 0;
	int lIndex;
	int rIndex;

	for (int i = 1; i < treeSize; i+=2){
    
	    if (i > 1) {
            parentIndex = i - levelIndex - ((numIndexesInLevel - levelIndex) / 2);
        }
		
        dpfPRG(ctx, seeds[parentIndex], &sL, &sR, &tL, &tR);

		if (bits[parentIndex] == 1){
			sL = sL ^ sCW[currLevel];
			sR = sR ^ sCW[currLevel];
			tL = tL ^ tCW0[currLevel];
			tR = tR ^ tCW1[currLevel];
		}

        lIndex =  i;
        rIndex =  i + 1;
        seeds[lIndex] = sL;
        bits[lIndex] = tL;
        seeds[rIndex] = sR;
        bits[rIndex] = tR;

        levelIndex += 2;
        if (levelIndex == numIndexesInLevel) {
            currLevel++;
            numIndexesInLevel *= 2;
            levelIndex = 0;
        }
    }

	uint128_t hashinput[hash1->outblocks];
	uint128_t tpi[hash1->outblocks];
	uint128_t cpi[hash2->outblocks];
	int index;
	int bit;

	for (int i = 0; i < numLeaves; i++) {
        index = treeSize - numLeaves + i;
      
		uint128_t res = convert(&seeds[index]);
		// printf("res[%i] = %lx \n", i, (unsigned long)res);

		// *********************************
		// START: DPF verification code 
		// *********************************
		bit = seed_lsb(seeds[index]);
		// printf("S%i l: %i bit %i ", b, l, bit);
		// printf(" pi = %llx", (unsigned long long)pi[0]);

        hashinput[0] = index;
		hashinput[1] = seeds[size];

		// step 1: H(seeds[size]||X[l])
		mmoHash2to4(hash1, (uint8_t*)&hashinput[0], (uint8_t*)&tpi[0]);

		// step 2: pi^correct(tpi, cs, bit)
		hashinput[0] = pi[0]^correct(tpi[0], cs[0], bit);
		hashinput[1] = pi[1]^correct(tpi[1], cs[1], bit);
		hashinput[2] = pi[2]^correct(tpi[2], cs[2], bit);
		hashinput[3] = pi[3]^correct(tpi[3], cs[3], bit);
		
		// step 3: H'(step2)
		mmoHash4to2(hash2, (uint8_t*)&hashinput[0], (uint8_t*)&cpi[0]);
		
		// pi = pi XOR H'(pi XOR correct(tpi, cs, bit))
		pi[0] ^= cpi[0]; 
		pi[1] ^= cpi[1];

		if (bit == 1) {
			//correction word
			res = res ^ convert((uint128_t*)&k[18*size+18]);
		}

		// copy block to byte output
		memcpy(&out[i*sizeof(uint128_t)], &res, sizeof(uint128_t));
    }

	free(bits);
	free(seeds);
}
