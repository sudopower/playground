#include<stdio.h>

/*
Fahrenheit - Celsius conversions in steps of 20 with 2 decimal step precision
*/

main(){
    float farh,celsius;
    int upper,lower;
    int steps=20;

    lower = 0;
    upper=300;

    farh=lower;
    printf("farenheit\tcelsius\n");
    while(farh<=upper){
        celsius=5*(farh-32)/9;
        printf("%3.0f\t\t%6.2f\n",farh,celsius);
        farh+=steps;
    }
}