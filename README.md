SMCFIX

<img alt="SMCFix icon, a SNES cart guillotine." src="https://github.com/fith/smcfix/blob/main/assets/icon.png?raw=true" width="128"/>

Cross-platform command line and GUI tool for removing headers on .smc files (SNES ROMs).

SMC headers are often added to hold metadata for software emulators, but
altering the original format breaks the ROM for FPGA emulation which expects
original-hardware accurate data.

Basically this is to fix games to play on my Analogue Pocket.

It's very fast and can process a whole directory or individual files. Overwrite
existing files or create new ones with a suffix.

Barebones right now. Might be a naive implementation, but has fixed all the broken
ROMs I found to test on.

Usage of ./smcfix:
  -dir string
    	Directory to scan for SMC files. (default "/Users/kevin/Workspace/smcfix/bin/mac")
  -file string
    	Single SMC file to check and clean.
  -help
    	Show this help.
  -out string
    	Specify alternate output directory.
  -overwrite
    	Overwrite or create new e.g. "[filename]-smcfix.smc" (default false)
      
For smcfix.app the command line utility would be run from smcfix.app/Contents/MacOS/smcfix
