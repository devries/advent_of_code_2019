CC=clang

.PHONY: all clean

part1: part1.o hashtable.o
	$(CC) -o $@ part1.o hashtable.o

part2: part2.o hashtable.o
	$(CC) -o $@ part2.o hashtable.o

%.o: %.c
	$(CC) -c -o $@ $<

all: part1 part2

clean:
	-rm part1.o part2.o part1 part2
