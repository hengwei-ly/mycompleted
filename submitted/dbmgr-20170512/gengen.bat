if "%gengen_path%" EQU "" (
  if exist D:\developing\go\meijing\tpt_vendor\src\github.com\runner-mei\gengen (
    set gengen_path=D:\developing\go\meijing\tpt_vendor\src\github.com\runner-mei\gengen
  )
  if exist D:\dev\tpt_vendor\src\github.com\runner-mei\gengen (
    set gengen_path=D:\developing\tpt_vendor\src\github.com\runner-mei\gengen
  )
)

%gengen_path%\gengen.exe %*
