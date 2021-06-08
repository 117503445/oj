#include <iostream>
using namespace std;

int main(int argc, char const *argv[])
{
    int a, b;
    cin >> a;
    cin >> b;

    string s = to_string(a + b);
    //cout << s;
    for (size_t i = 0; i < s.length(); i++)
    {
        cout << s[i];
        if ((s.length() - i) % 3 == 1 && (s.length() - i) != 1 && s[i] != '-')
        {
            cout << ',';
        }
    }

    return 0;
}
