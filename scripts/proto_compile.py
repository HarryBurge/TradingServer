import subprocess
import os

go_module = "TradingServer"

targets = [
    ("internal/api/grpc/dataset/*.proto", "internal/api/grpc/dataset/pb"),
    ("internal/api/grpc/client/*.proto", "internal/api/grpc/client/pb")
]

for (in_file, out_file) in targets:
    if not os.path.isdir(out_file):
        os.mkdir(out_file)

    subprocess.run(
        f"protoc --go_out=./{out_file} --go_opt=module={go_module}/{out_file} --go-grpc_out=./{out_file} --go-grpc_opt=module={go_module}/{out_file} .\{in_file}",
        shell=False
    )