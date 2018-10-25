package getchar

//#include <stdio.h>
//#include <termios.h>
//#include <unistd.h>
// int getch(void)
// {
//     struct termios oldattr, newattr;
//     int ch;
//     tcgetattr( STDIN_FILENO, &oldattr );
//     newattr = oldattr;
//     newattr.c_lflag &= ~( ICANON | ECHO );
//     tcsetattr( STDIN_FILENO, TCSANOW, &newattr );
//     ch = getchar();
//     tcsetattr( STDIN_FILENO, TCSANOW, &oldattr );
//     return ch;
// }
import "C"

func Rune() rune {
	runeInt := C.getch()
	return rune(runeInt)
}
