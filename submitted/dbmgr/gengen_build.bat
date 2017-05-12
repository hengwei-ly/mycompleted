
if "%gengen_path%" EQU "" (
  if exist D:\developing\go\meijing\tpt_vendor\src\github.com\runner-mei\gengen (
    set gengen_path=D:\developing\go\meijing\tpt_vendor\src\github.com\runner-mei\gengen
  )
  if exist D:\dev\tpt_vendor\src\github.com\runner-mei\gengen (
    set gengen_path=D:\dev\tpt_vendor\src\github.com\runner-mei\gengen
  )
)

@pushd %gengen_path%
@%gengen_path%\gengen.exe embeded template_gen.go
@goimports  -w template_gen.go
@go build
@%gengen_path%\gengen.exe embeded template_gen.go
@goimports  -w template_gen.go
@go build
@popd
