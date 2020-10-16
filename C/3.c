#include<stdio.h>
main(){
    const digits=30;
    int dec,count;
    int arr[digits];
    printf("Enter decimal:");
    scanf("%d",&dec);
    while(dec>0){
        arr[count]=dec%2;
        dec=(int)dec/2;
        count++;
    }
    printf("binary:");
    do{
        count--;
        printf("%d",arr[count]);
    }while(count>0);
}