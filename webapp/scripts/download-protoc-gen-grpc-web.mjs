import * as os from 'node:os';
import * as fs from 'node:fs';

function platform() {
  switch (os.platform()) {
    case 'darwin':
      return 'darwin';
    case 'linux':
      return 'linux';
    default:
      throw new Error(`Unsupported platform: ${os.platform()}`);
  }
}

function arch() {
  switch (os.arch()) {
    case 'x64':
      return 'x86_64';
    case 'arm':
    case 'arm64':
      return 'aarch64';
    default:
      throw new Error(`Unsupported architecture: ${os.arch()}`);
  }
}

function version() {
  return '1.5.0';
}

function url() {
  return `https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-${version()}-${platform()}-${arch()}`;
}

const targetDir = './node_modules/.bin';
const downloadUrl = url();

console.log('Downloading protoc-gen-grpc-web from', downloadUrl);
fetch(downloadUrl).then((res) => {

  if (res.status !== 200) {
    throw new Error(`Failed to download protoc-gen-grpc-web: ${res.status}`);
  }

  return res.arrayBuffer();
}).then((data) => {
  fs.writeFileSync(targetDir + '/protoc-gen-grpc-web', Buffer.from(data));
  fs.chmodSync(targetDir + '/protoc-gen-grpc-web', 0o755);
});