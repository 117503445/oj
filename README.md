# OJ

Automatically compile and test command line programs on the fly.

## About The Project

When solving algorithmic problems, you need to continuously compile and input test cases to determine whether the output results are correct. This product can help you automate these processes, you only need to write code.

## Getting Started

### Installation

1. According to your own system, download the latest binary file from <https://github.com/117503445/oj/releases/latest> and add folder to the system `path` environment variable.

2. input `oj` in terminal, show `No Source Code found in the dir.`.

3. According to your needs, install gcc, g++, python, and make sure it is accessible in the terminal.

### Usage

We recommend that you use vscode and enable auto-save.

1. We assume that you write cpp code in ~/mydir/main.cpp.

    ```cpp
    #include <iostream>
    using namespace std;

    int main(int argc, char const *argv[])
    {
        cout << "start\n";

        int a, b;
        cin >> a >> b;
        cout << "a = " << a << endl;
        cout << "b = " << b << endl;
        
        cout << a + b;
        cout << "\nend";
        return 0;
    }
    ```

2. write input test file.

    ~/mydir/1.in

    ```text
    1 2
    ```

    ~/mydir/2.in

    ```text
    2 3
    ```

3. run `oj` in `~/mydir`

4. When `main.cpp` changes, `oj` will automatically compile `main.cpp`, redirect `1.in` and `2.in` as test cases from standard input to a binary file, and save the output to `1.main.out` and `2.main.out` .

### Dev

install golang v1.22

```sh
go run /workspace/cmd/oj
```

There are some test code in ./assets dir.

## License

Distributed under The GNU General Public License v3.0, See `LICENSE` for more information.
