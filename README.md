# Advent of Code 2019

## Summary of My Experience

Last year was when I first heard of [Advent of
Code](https://adventofcode.com/2019), but I was studying for some
certifications and starting a new job and only completed the first two days.
This year, I was reminded in November and decided to go for it. My goal was to
complete each exercise by the end of the day (before work if possible), though
I was not planning to get up at midnight (when the puzzle is released in my
time zone). I managed to complete each puzzle by the end of the day except for
Day 14 when I was traveling from Washington DC to Maine and decided to get
some sleep instead. I also convinced a brilliant friend to do Advent of Code
with me. We shared a leaderboard, and keeping up with him added some extra
pressure to get all the stars.

I chose to use Go because it is a language which I have started using a bit in
the last year, and there is no better way to jump in than doing a lot of
programming. In several cases it might have been easier for me to use Python
as I was already familiar with everything I needed in Python, whereas with Go
I had to look up a number of things. I found Go to be a very capable language
for the competition. I took advantage of the default map return values often,
and used the map for any place where I would have used a set in python. It was
easy to use bitfields in Go, which came in handy on Day 18 and Day 24. I
decided to be explicit about using `int64` in my Intcode, and that led to a
lot of casting `int` values to `int64` which was a little tedious. I also
tried to make use of unit tests, so I became more familiar with the `go test`
command than I had been.

I really like the Intcode problems, and I have marked all the problems using
Intcode in my index below. Each directory shows the intcode machine as it was
developed and the various changes as we went up through the last day. It was
relatively stable after Day 9, but I did make a couple of small modification.
I took advantage of Go's goroutines and channels to manage input and output.
This also allowed me to run several interpreters simultaneously. The only bump
I had with Intcode was on Day 23 when I had trouble figuring out when all the
intcode computers were awaiting input. I think I could have used a
`sync.WaitGroup` for that, but I decided to just have a 1 second timeout in
the router which sent the 255 packet to 0 if no other traffic was seen on the
network. Working on Intcode reminded me a lot of working on problems in *The
Elements of Computing Systems: Building a Modern Computer from First
Principles* by Nisan and Schocken, which I would recommend to anyone who
enjoyed Intcode.

I managed to work without any hints except for Day 12, when I just couldn't
figure out how to do part 2. It didn't occur to me to find how long it took
for each dimension to cycle and then use the least common multiple of those
cycles. I looked for hints on reddit when I felt like I couldn't think of any
way to proceed.

Day 13 was one of my favorites, and one of the only ones for which I made [a
visualization](https://youtu.be/pWBlPCahKQw). At first I was going to predict
the path and move to the position where I would next need to be to reflect the
ball, but in the end I found it easier to just remain under the ball. I also
liked Day 15 a lot.

Day 18 was one of the hardest for me. I managed to do part A with a breadth
first search, but then in part B I just had too large a parameter space. I
decided to think about the problem more before looking for help, and I started
going down the path of solving each quadrant separately. I found that each
quadrant was easily done with a BFS, so I decided to see if I could send each
robot to the point where it needed a key from another quadrant before
switching to the next robot. I managed to prove to myself that if the puzzle
was possible, that a robot would just have to get the keys in its own quadrant
in the right order and have the optimal path, so I just had to add up the
steps for all four quadrants and it worked. I finished that puzzle at
23:59:13.

I decided to jump into Day 19 since I was already up submitting day 18, and
for some reason part B just clicked for me, and I was done in only 30 minutes
and managed to get close to the leaderboard (117th) which is not something I
expected to pull off.

I think my best solution was Day 22 part B. I was working to see if I could
combine my shuffle instructions into one simple function, when I realized cut
and deal with increment were linear functions. I checked that deal into new
stack was also linear, and then began making a struct for coefficients of
linear functions. Pretty soon I had a nice way to consolidate my shuffle
procedure into a single linear function, and I started working on how to
repeat that function. I saw that the constant was a geometric series, and then
had to look up how to do modular division. Luckily Go's big integer library
had both a modular exponential formula and a modular inverse, so I was able to
implement both the repeated linear function as well as the inverse linear
function as required by part B. I learned a lot about modular arithmetic in
doing that problem, and was very glad I was able to figure it out. I never
would have been able to do it without writing out some equations on paper, so
I am in awe of people who worked it out just while typing on the computer.

On December 21st I saw [a
post](https://www.reddit.com/r/adventofcode/comments/edl79n/intcode_textbased_adventure/)
by `sbguest` about a text adventure he or she had written in Intcode. I
decided to put together a version of my Intcode machine I could use via the
keyboard to play that game. This code worked unchanged for December 25th.
Finally all my experience playing MUDs was paying off, and I managed to be the
13th to finish using only manual entry.

## Index

- [Day 1: The Tyranny of the Rocket
  Equation](https://adventofcode.com/2019/day/1) [part 1](day01_p1) [part
  2](day01_p2)
- [Day 2: 1202 Program Alarm](https://adventofcode.com/2019/day/2) (Intcode) [part 1](day02_p1) [part 2](day02_p2)
- [Day 3: Crossed Wires](https://adventofcode.com/2019/day/3) [part 1](day03_p1) [part 2](day03_p2)
- [Day 4: Secure Container](https://adventofcode.com/2019/day/4) [part 1](day04_p1) [part 2](day04_p2)
- [Day 5: Sunny with a Chance of
  Asteroids](https://adventofcode.com/2019/day/5) (Intcode) [part 1](day05_p1) [part
  2](day05_p2)
- [Day 6: Universal Orbit Map](https://adventofcode.com/2019/day/6) [part 1](day06_p1) [part 2](day06_p2)
- [Day 7: Amplification Circuit](https://adventofcode.com/2019/day/7) (Intcode) [part 1](day07_p1) [part 2](day07_p2)
- [Day 8: Space Image Format](https://adventofcode.com/2019/day/8) [part 1](day08_p1) [part 2](day08_p2)
- [Day 9: Sensor Boost](https://adventofcode.com/2019/day/9) (Intcode) [part 1](day09_p1) [part 2](day09_p2)
- [Day 10: Monitoring Station](https://adventofcode.com/2019/day/10) [part 1](day10_p1) [part 2](day10_p2)
- [Day 11: Space Police](https://adventofcode.com/2019/day/11) (Intcode) [part 1](day11_p1) [part 2](day11_p2)
- [Day 12: The N-Body Problem](https://adventofcode.com/2019/day/12) [part 1](day12_p1) [part 2](day12_p2)
- [Day 13: Care Package](https://adventofcode.com/2019/day/13) (Intcode) [part 1](day13_p1) [part 2](day13_p2)
- [Day 14: Space Stoichiometry](https://adventofcode.com/2019/day/14) [part 1](day14_p1) [part2](day14_p2)
- [Day 15: Oxygen System](https://adventofcode.com/2019/day/15) (Intcode) [part 1](day15_p1) [part 2](day15_p2)
- [Day 16: Flawed Frequency
  Transmission](https://adventofcode.com/2019/day/16) [part 1](day16_p1) [part 2](day16_p2)
- [Day 17: Set and Forget](https://adventofcode.com/2019/day/17) (Intcode) [part 1](day17_p1) [part 2](day17_p2)
- [Day 18: Many-Worlds Interpretation](https://adventofcode.com/2019/day/18) [part 1](day18_p1) [part 2](day18_p2)
- [Day 19: Tractor Beam](https://adventofcode.com/2019/day/19) (Intcode) [part 1](day19_p1) [part 2](day19_p2)
- [Day 20: Donut Maze](https://adventofcode.com/2019/day/20) [part
  1](day20_p1) [part 2](day20_p2)
- [Day 21: Springdroid Adventure](https://adventofcode.com/2019/day/21)
  (Intcode) [part 1](day21_p1) [part 2](day21_p2)
- [Day 22: Slam Shuffle](https://adventofcode.com/2019/day/22) [part
  1](day22_p1) [part 2](day22_p2)
- [Day 23: Category Six](https://adventofcode.com/2019/day/23) (Intcode) [part
  1](day23_p1) [part 2](day23_p2)
- [Day 24: Planet of Discord](https://adventofcode.com/2019/day/24) [part
  1](day24_p1) [part 2](day24_p2)
- [Day 25: Cryostasis](https://adventofcode.com/2019/day/25) (Intcode) [part
  1](day25_p1)
