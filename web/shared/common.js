// Shared WASM loader utility for claudefun web apps
async function loadWasm(wasmPath) {
  if (!WebAssembly) {
    throw new Error('WebAssembly not supported in this browser');
  }
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(wasmPath),
    go.importObject
  );
  go.run(result.instance);
  return result.instance;
}
