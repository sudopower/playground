fn main() {
    println!("Hello world");
    println!("I'm Kiran !");

    println!("I can print using  println!() and I need a semicolon at end!");

    let nos = 1;

    println!("I can print a number example - {}", nos);

    print!("I can print on the same line using");
    print!(" print()! \n");

    println!("I can format the code using cargo fmt");

    println!("I can also format it using rustfmt filename");

    let sum = nos + 1;
    println!("I can add nos using num1 + num2 = {}", sum);

	printAGreeting();

	printANumber(67);
}

fn printAGreeting(){
	println!("This greeting was printed from a function")
}

fn printANumber(x i64){
	println!("{}: This number was printed from a arg to a function",x)
}
