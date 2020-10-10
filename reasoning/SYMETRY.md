# Coping With Symetry


We represent a square dance formation using a struct.  In the case of
a two dancer formation, that struct will have two slots: one for each
dancer.  The slots of a struct are necessarily distinct.

In a formation like Couple or Tandem, the two dancers have distinct
roles: beau/belle and leader/trailer respectively.

Other formations like MiniWave, FaceToFace or BackToBack are
symmetric.  This can lead to redundancy in our representations.  For
example:

<pre>
&MiniWaveImpl(dancer_1, dancer_2)
</pre>

and

<pre>
&MiniWaveImpl(dancer_2, dancer_1)
</pre>

are two different representations of the same MiniWave.

The rule for recognizing MiniWaves could avoid emitting redundant
MiniWaves by imposing a relationship between the Ordinals of the two
dancers.  For example

>pre>
    if dancer1.Ordinal() >= dancer2.Ordinal() {
      	return
    }
</pre>

The rules for recognizing larger formations typically take smaller
formations as input and check that the dancer with a given role in one
input formation has another role in another input formation.

But formations that derive from MiniWave (or any other symetric two
dancer formation) don't have distinct roles to compare on, so rules
resort to using HasDancer rather than formation role member equality.


How to detect a WaveOfFour given its three component MiniWaves?

1. Center MiniWave must have opposite handedness from the other two.

2. Other two must be different from each other.

3. Each "end" MiniWave must share a dancer with the center MiniWave.

Generalize HasDancer to an intersection: CommonDancers.  Thus, the
test for the third condition might be

<pre>
   len(center.CommonDancers(side1)) == 1
 </pre>

We might make it a function rather than a method so that we can use
interface types for both arguments:

<pre>
func CommonDancers(d1, d2 Dancers) []Dancers {
    ...
}
</pre>

Or maybe just name it Intersection.  We already have dancer.Intersection.

The center MiniWave has ordered dancers.  Do we want to impose that
the first side MiniWave contains the first dancer of the center
MiniWave?

Do we want to represent the WaveOfFour as the three MiniWaves or as
Dancer, Dancer, Dancer, and Dancer?  This representation looses
identity of the MiniWaves.

Perhaps a measure of the utility of a formation is how easy it is to
figure out how to do "centers cross run" or "cross circulate".
Regular circulate isn't an issue because we have Centers and Ends.
What about circulate from a TidalWave?  Do we need a notion of "nth
from each end" or nth from center?

