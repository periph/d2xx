# d2xx driver

All files in this directory are from the Future Technology "D2XX" driver at
https://ftdichip.com/drivers/d2xx-drivers/.

## License

> This software is provided by Future Technology Devices International Limited
> ``as is'' and any express or implied warranties, including, but not limited
> to, the implied warranties of merchantability and fitness for a particular
> purpose are disclaimed. In no event shall future technology devices
> international limited be liable for any direct, indirect, incidental, special,
> exemplary, or consequential damages (including, but not limited to,
> procurement of substitute goods or services; loss of use, data, or profits; or
> business interruption) however caused and on any theory of liability, whether
> in contract, strict liability, or tort (including negligence or otherwise)
> arising in any way out of the use of this software, even if advised of the
> possibility of such damage.  FTDI drivers may be used only in conjunction with
> products based on FTDI parts.
>
> FTDI drivers may be distributed in any form as long as license information is
> not modified.
>
> If a custom vendor ID and/or product ID or description string are used, it is
> the responsibility of the product manufacturer to maintain any changes and
> subsequent WHCK re-certification as a result of making these changes.
>
> For more detail on FTDI Chip Driver licence terms, please [click
> here](http://www.ftdichip.com/Drivers/FTDriverLicenceTermsSummary.htm).


## Modifications

- Fixed ftd2xx.h to UTF-8
- Added a minimal copyright header from WinTypes.h
- Converted header files from CRLF to LF
- Removed trailing spaces and tabs
- Renamed the static libraries to include the platform and the version.
