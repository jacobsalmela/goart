<p align="center">
    <a href="https://github.com/jacobsalmela/goart">
        <img src="https://user-images.githubusercontent.com/3843505/138571722-766c574f-78ed-4c98-88a9-328e95ba6b53.png" width="450" height="250" alt="goart ">
    </a>
    <br>
    <strong>goart</strong><br>
    Generate art on the command line.
</p>

This command-line utility, `goart`, began as a fork of the companion source code to the book, ["Generative Art in Go."](https://preslav.me/generative-art-in-golang/), which intended to introduce novice and experienced programmers to algorithmic art.  It did just that for me, and I decided to convert the code into a command-line utility and being modifying the code to do even more.

```
$ goart --help
Generates art using a source image as the starting point.

        The resulting image can be further manipulated with flags to this
        command.

Usage:
  goart generate FILE [flags]

Flags:
  -A, --alpha-increase float               the step of transparency increase at each iteration (default 0.06)
  -H, --height int                         height of the generated image (default 2000)
  -h, --help                               help for generate
  -a, --initial-alpha float                beginning stroke transparency (default 0.1)
  -K, --keep-source-dimensions             generate a new image with the same dimensions as the source image (default true)
  -M, --max-edge int                       maximum stroke will be a n-edge polygon (4 is a square) (default 4)
  -m, --min-edge int                       minimum stroke will be a n-edge polygon (3 is a triangle) (default 3)
  -z, --random                             generates a new image using a random one pulled from source.unsplash.com
  -t, --stroke-inversion-threshold float   the minimum stroke size (default 0.05)
  -j, --stroke-jitter float                deviation of the colored stroke from its projected position in the original image (default 0.1)
  -r, --stroke-ratio float                 size of the initial stroke compared to that of the final result (default 0.75)
  -R, --stroke-reduction float             the initial stroke size gets minimized by this amount on each iteration (default 0.002)
  -T, --total-cycles int                   Copy any discovered k8s squashfs images from SRC to DEST (default 5000)
  -W, --width int                          width of the generated image (default 2000)

Global Flags:
      --config string   config file (default is $HOME/.goart.yaml)
      --viper           use Viper for configuration (default true)
```