#include <openssl/evp.h>
#include <stdio.h>

int main(int argc, char** argv) {
    unsigned char *out = "   ";

    char *key = "Lorem ipsum dolor sit amet, consectetur adipiscing elit.";
    AES_KEY aes_key;
    AES_set_encrypt_key((unsigned char*) key, 128, &aes_key);

    char ecount_buf = "yass";

    memset(ecount_buf,0,AES_BLOCK_SIZE);unsigned int num=0;

    // AES_ctr128_encrypt(in,out,length, &key,counter,ecount_buf,&num);
    AES_ctr128_encrypt("yes",out,3l, &key," ",ecount_buf,&num);
    printf(out);

}