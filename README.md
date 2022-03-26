# Volume

Thanks for the task. I enjoyed solving it!

## What is the tricky part hidden in the challenge?
It becomes obvious quickly that in this case it's all about forks and loops. In the description
you gave an examples of linear journeys only.

## Is every set of hops valid?
Obviously not. There might be forks leading to nowhere, loops which make detecting starting and finishing points impossible,
or set of hops may be simply empty. I included unit tests describing different cases in [./reduce_test.go](./reduce_test.go).

## Do I have to find the exact path?
No! I'm interested only in starting and finishing points. Detecting full path has complexity O(n^2), while detecting first
and last points takes only O(n).

## How to describe the problem?
Airports are vertices of the graph, hops are its directed edges.

## How to detect if graph is solvable?
Conditions graph has to meet to be solvable:
- there is exactly one vertex having exactly one more outgoing edges than incoming ones - this is starting point
- there is exactly one vertex having exactly one more incoming edges than outgoing ones - this is finishing point
- there are 0 or more (but finite) vertices having exactly the same number of incoming and outgoing edges

## So what am I dealing with exactly?
After looking at rules mentioned in the previous paragraph it becomes clear that I'm dealing with acyclic Eulerian graph and
the task is about finding starting and finishing points of Eulerian path inside that graph.

Because those points may be detected just by looking at the edges around these single points I don't need to traverse the graph.

## Algorithm
1. Count incoming and outgoing edges around each vertex of graph formed from hops. 
2. Verify (by looking at counted numbers) that graph meets rules of acyclic Eulerian graph.
3. Select starting and finishing points.

## Loop
There is one interesting case. Suppose this is the list of hops:

`(AAA, BBB), (BBB, CCC), (CCC, AAA)`

Graph built from those edges forms a correct trip, however it is an *cyclic* Eulerian graph instead of acyclic one.
So despite the fact that the trip exists, I'm not able to select starting and finishing points.

Making a statement that graph has to be acyclic I implicitly treat this case as invalid.