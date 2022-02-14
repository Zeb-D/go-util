//oc2go.c
#include <stdio.h>
#include <string.h>
#include "oc2go.h"                       //此处为上一步生成的.h文件

int main(){
    char c1[] = "did";
    GoString s1 = {c1,strlen(c1)};       //构建Go语言的字符串类型
    char *c = helloV5(s1);
    printf("r:%s",c);
    return 0;
}