import os
import platform
import subprocess
import shutil

PROJECT_NAME = "MCOneBot"
VERSION = "1.0.0"
MAIN_FILE = "main.go"
PLATFORMS = ["windows", "linux", "darwin"]
ARCHS = ["amd64", "386", "arm64", "arm"]
OUTPUT_PATH = "build"

def compile_for_platform(system, arch):
    output_name = f"{PROJECT_NAME}-{VERSION}-{system}-{arch}"
    output_file = os.path.join(OUTPUT_PATH, output_name)
    if system == "windows":
        output_file += ".exe"

    env = os.environ.copy()
    env["GOOS"] = system
    env["GOARCH"] = arch

    command = ["go", "build", "-o", output_file, MAIN_FILE]

    print(f"Compiling for {system} {arch}...")
    subprocess.run(command, env=env)

def main():
    if os.path.exists(OUTPUT_PATH):
        shutil.rmtree(OUTPUT_PATH)
    os.makedirs(OUTPUT_PATH)

    for system in PLATFORMS:
        for arch in ARCHS:
            compile_for_platform(system, arch)

    # Compile for the current platform
    current_system = platform.system().lower()
    current_arch = platform.machine().lower()
    output_name = f"{PROJECT_NAME}-{VERSION}-{current_system}-{current_arch}"
    output_file = os.path.join(OUTPUT_PATH, output_name)
    if current_system == "windows":
        output_file += ".exe"
    print(f"Compiling for {current_system} {current_arch}...")
    subprocess.run(["go", "build", "-o", output_file, MAIN_FILE])

if __name__ == '__main__':
    main()
