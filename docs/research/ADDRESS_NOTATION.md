# Address Spaces

Always prefix an address with its domain:

- `ROMCPU:$C10E14`
- `ROMFILE:0x010E14`
- `WRAM:$7E2EB5`
- `WRAM:+$2EB5`
- `VRAM:$4000`
- `CGRAM:$0080`
- `OAM:$0120`
- `ARAM:$3000`
- `DSP:$5D`
- `SRAM:+$0000` — cartridge save RAM, offset from the start of the save
  region. Added 2026-08-01 when `ff6lab state` began reading the
  `cart.saveRam` block out of savestates; no FF6 SRAM layout is
  established yet.

A naked hexadecimal number is not acceptable in canonical documentation.
