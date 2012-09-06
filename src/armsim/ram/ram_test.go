package ram

import "testing"

func TestInit(t *testing.T) {
  // Test the RAM initializer at 32k
  runInit(t, 32 * 1024)
  // Test the RAM initializer at 32m
  runInit(t, 32 * 1024 * 1024)
}

// Helper method for the RAM initializer testing code
func runInit(t *testing.T, nBytes uint32) {
  // Setup
  ram := Init(nBytes)

  // Ensure RAM is the same length as nBytes
  if (uint32(len(ram.memory)) != nBytes) {
    t.Fatal("Did not allocate RAM to specified nBytes")
  }

  // Ensure RAM is initialized to zero
  for i := uint32(0); i < nBytes; i++ {
    if ram.memory[i] != byte(0) {
      t.Fatal("Did not initalize RAM sector", i)
    }
  }
}

func TestWriteByte(t *testing.T) {
  // Setup
  ram := Init(32)
  var byteToWrite byte = 95;

  // Test under normal conditions
  result := ram.WriteByte(0, byteToWrite)
  if result == true && byteToWrite != ram.memory[0] {
    t.Fatal("Did not properly write byte.")
  }

  // Test address out of range
  result = ram.WriteByte(33, byteToWrite)
  if result == true {
    t.Fatal("Attempted to write out of range address.")
  }
}

func TestReadByte(t *testing.T) {
  // Setup
  ram := Init(32)
  var dataToWrite byte = 0xFF

  // Test normal reading of bytes
  ram.WriteByte(0, dataToWrite)
  b, result := ram.ReadByte(0)
  if !result || b != 0xFF {
    t.Fatalf("expected %#x got %#x", dataToWrite, b)
  }

  // Test address out of range
  b, result = ram.ReadByte(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}

func TestWriteHalfWord(t *testing.T) {
  // Setup
  ram := Init(32)

  // Test normal writing of halfwords
  // Lots of tests because this is faily complicated with shifts and casts.
  var hwToWrite uint16

  hwToWrite = 0xFFFF
  result := ram.WriteHalfWord(0, hwToWrite)
  b1, _ := ram.ReadByte(0)
  b2, _ := ram.ReadByte(1)

  if !result || b1 != 0xFF || b2 != 0xFF {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0xFF00
  result = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if !result || b1 != 0xFF || b2 != 0x00 {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  hwToWrite = 0x0F0D
  result = ram.WriteHalfWord(0, hwToWrite)
  b1, _ = ram.ReadByte(0)
  b2, _ = ram.ReadByte(1)

  if !result || b1 != 0x0F || b2 != 0x0D {
    t.Fatal("Did not properly write half word. Expected:", hwToWrite, "Got:", b1, b2)
  }

  // Test writing to odd address
  result = ram.WriteHalfWord(1, hwToWrite)

  if result {
    t.Fatal("Wrote halfword to odd address. This should not happen.")
  }

  // Test writing to out-of-bounds address
  result = ram.WriteHalfWord(34, hwToWrite)
  if result {
    t.Fatal("Wrote halfword to an out-of-bounds address.")
  }
}

func TestReadHalfWord(t *testing.T) {
  // Setup
  ram := Init(32)
  var hwToWrite uint16 = 0xFFFF

  // Test normal reading of bytes
  ram.WriteHalfWord(0, hwToWrite)
  b, result := ram.ReadHalfWord(0)
  if !result || b != hwToWrite {
    t.Fatalf("expected %#x got %#x", hwToWrite, b)
  }

  // Test address out of range
  _, result = ram.ReadHalfWord(80)
  if result {
    t.Fatal("Attempted to read out of range address.")
  }
}
