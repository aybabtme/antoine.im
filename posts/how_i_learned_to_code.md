{
    "title":"An Infantryman Learns To Code",
    "author":"Antoine Grondin",
    "date":"2016-02-25T23:59:59.000Z",
    "invisible": false,
    "abstract":"In 2009-2010, I was an infantry officer working in Kandahar (Afghanistan) Airfield's brigade HQ; a Canadian HQ then.",
    "language":"en"
}

Before getting into this software thing, I was an infantry officer. My job
description was something like:

> An Infantry Officer performs a wide range of duties, from commanding and
> leading soldiers [...] operating anywhere in the world, in any environment
> [...]. The primary role of Infantry during operations is to be involved in
> combat.

On the job, I looked something like:

<a href="/assets/data/how_i_learned_to_code/back_in_infantry.jpg">
    <img src="/assets/data/how_i_learned_to_code/back_in_infantry.jpg" style="width:60%;"/>
</a>

Somehow I ended up learning how to code while on tour in Afghanistan.

## Fast backward to 2009

In 2009-2010, I was an infantry officer working in [Kandahar][]
Airfield for a brigade HQ known as [Task Force Kandahar][TFK]; a Canadian HQ
then. My job was to receive all the reports from the ground and fan that info
out throughout the HQ and our higher ups/sides.

The HQ was processing things like medical evacuations, support fire missions,
contact reports. The way it was done was very frustratingly inefficient. For
instance, a medical mission would go like this:

* A unit would report an IED strike with critical injuries (using a [9-liner][]).
* The unit would pass a [MEDEVAC][] request by radio to
their company -> battalion -> brigade HQ (us).
* We would then be the dispatch center for the helicopters and synchronization of airspace and hospital and all.

The request would arrive ~30s to 1m after the actual strike. I would yell at
an airman that would get out of his chair, walk to the center map and with
his ruler, measure the distance in miles between the
hospital landing pad and the strike. He would then compute ETAs for the
helicopters based on various parameters, and then ask the helicopter HQ to
send a chopper on site. Only then would he slowly type a message in a predefined
format, and then post that text in the communication channel (that's ~5m after
the strike).

This 5m latency in sending the request to choppers, and giving back ETD/ETA
info to the unit on the ground, would result in people dying from their wounds
or staying in dangerous/exposed locations longer than strictly necessary
(waiting for chopper ETAs). This 5m latency was putting people at risk and
killing folks.

The whole thing was very frustrating. The man's job was crying for computer
automation, but I didn't know anything about programming at the time. I googled
a bit and came up with people saying that to program, you had to learn Java or
C, and needed a compiler for them. I didn't know anything about this whole
"programming" thing back then, the word didn't mean much to me. You can only
imagine what the officer of the Signal Squadron thought when I asked:

> "Can you install a C compiler on my computer?"

## Would you give a C compiler to an Infantry Officer?

Would you? The Signals Officer sort of laughed at me and told me to use Excel.
So I wrote a tool to automate this man's job in Excel. I picked up what I could
using fancy techniques like the Pythagorean theorem and string concatenation.

<a href="/assets/data/how_i_learned_to_code/maintainable_code.png">
    <img src="/assets/data/how_i_learned_to_code/maintainable_code.png" style="width:60%;"/>
</a>

In the end, the tool was very crude but accomplished something very
useful:

- It had a flow that ensured all the reports required by people on the ground,
and above, were sent in a timely and orderly manner.
- Each step of that flow was almost entirely automated.
  1. Each button filled a template and put the text in the clipboard for
  copy-pasting in the chat.
  2. Events were timed automatically.
  3. Distances and time of travel were computed automatically.
  4. A dropdown menu facilitated entering common values.
- Big warning signs were visible when a time critical step was ongoing, or some
important data was missing.

The result was this spreadsheet:

<a href="/assets/data/how_i_learned_to_code/overview.png">
    <img src="/assets/data/how_i_learned_to_code/overview.png" style="width:60%;"/>
</a>

... and it reduced the latency of "_receiving 9 liner to returning an ETA_" from
 ~5m to ~15s.

Being the Canadian army, they don't have this tendency our southern allies have
to give medals for eating breakfast. But I still got this:

<a href="/assets/data/how_i_learned_to_code/air_wing_comd_commandation_lt_grondin_aug_2010.jpg">
    <img src="/assets/data/how_i_learned_to_code/air_wing_comd_commandation_lt_grondin_aug_2010.jpg" style="width:60%;"/>
</a>

Eventually I figured out that there was a scripting language behind Excel
called [VBA][]. So I wrote tools to automate many other parts of the HQ, like a
tool to manage airspace for fire missions or a database to handle multiple
concurrent critical incidents:

<a href="/assets/data/how_i_learned_to_code/incident_manager.png">
    <img src="/assets/data/how_i_learned_to_code/incident_manager.png" style="width:60%;"/>
</a>

<a href="/assets/data/how_i_learned_to_code/9liner.png">
    <img src="/assets/data/how_i_learned_to_code/9liner.png" style="width:60%;"/>
</a>

After writing all that code, someone noticed I guess. Our operations were
running much better, the ability of our [TOC][] to handle multiple concurrent
incidents was greatly increased. So I also got a coin from Canada's Chief of
 Defence Staff (CDS, the top general at the time):

<a href="/assets/data/how_i_learned_to_code/important_people.jpg">
    <img src="/assets/data/how_i_learned_to_code/important_people.jpg" style="width:60%;"/>
</a>

> From left to right: unknown, [General Walter Natynczyk][Natynczyk] (used to be CDS), myself, [General Jonathan Vance][Vance] (current CDS), [Peter MacKay][MacKay] (used to be Minister of National Defence).

Then I realized I liked this programming thing much more than running around
with guns.

Then I came back to Canada (2010) and started a degree in software
engineering (2011).

Then I learned to code in [Go][] (2012) and interned at
Amazon (2013).

Then I interned at [Shopify][] (2014).

Then [Ben][] and
[Moisey Uretsky][Moisey] made a bet on me and offered me a [job][] at
[DigitalOcean][] (2014).

Then I graduated (2016).

## Then I'm here today.

When I joined the army, I felt I had begun a second life. A life where I would
learn to be a disciplined and organized person. When I started my degree, I
began a third life. Now I graduated, I'm not sure what to do. It's the first
time in years that I have so much time on my hands, with only a full time job.
Luckily for me, DigitalOcean is young enough a startup that I can have as large
an impact as I want by just _doing_ stuff.

I can recall when I started writing that VBA stuff, as I realized that I was
doing something that mattered, that I was changing the status quo. And I
realized it's not hard to do different, to do something that matters. All you
need to do is overcome inertia and get started.

## Note

Here's a sample of the first lines of code I ever wrote.  (warning: it may hurts
   your feelings)

https://gist.github.com/aybabtme/89e1f475c731ae985a1f


[Kandahar]: https://en.wikipedia.org/wiki/Kandahar_International_Airport
[TFK]: https://en.wikipedia.org/wiki/Task_Force_Kandahar
[9-liner]: http://www.army.mil/e2/c/downloads/355651.pdf
[MEDEVAC]: https://en.wikipedia.org/wiki/Medical_evacuation
[DUSTOFF]: https://en.wikipedia.org/wiki/Casualty_evacuation
[Natynczyk]: https://en.wikipedia.org/wiki/Walter_Natynczyk
[Vance]: https://en.wikipedia.org/wiki/Jonathan_Vance
[MacKay]: https://en.wikipedia.org/wiki/Peter_MacKay
[VBA]: https://en.wikipedia.org/wiki/Visual_Basic_for_Applications
[TOC]: https://en.wikipedia.org/wiki/Tactical_operations_center
[Go]: https://golang.org/
[Shopify]: https://www.shopify.com/
[Ben]: https://twitter.com/benuretsky
[Moisey]: https://twitter.com/moiseyuretsky
[DigitalOcean]: https://www.digitalocean.com/
[job]: http://grnh.se/wv3fgo
