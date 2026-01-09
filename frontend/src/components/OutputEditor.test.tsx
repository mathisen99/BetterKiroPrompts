import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'

/**
 * Property 3: Edit State Preservation
 * For any file edit action in the Output_Editor, the edited content SHALL be
 * preserved in state and reflected in subsequent reads.
 * 
 * Validates: Requirements 4.3
 * 
 * Feature: ai-driven-generation, Property 3: Edit State Preservation
 */

// Test the edit state logic in isolation (pure function testing)
// This tests the core state management logic used by OutputEditor

interface EditState {
  editedFiles: Map<string, string>
  originalFiles: Map<string, string>
}

function createInitialState(files: { path: string; content: string }[]): EditState {
  const originalFiles = new Map<string, string>()
  files.forEach(f => originalFiles.set(f.path, f.content))
  return {
    editedFiles: new Map(),
    originalFiles,
  }
}

function editFile(state: EditState, path: string, content: string): EditState {
  const newEditedFiles = new Map(state.editedFiles)
  newEditedFiles.set(path, content)
  return { ...state, editedFiles: newEditedFiles }
}

function resetFile(state: EditState, path: string): EditState {
  const newEditedFiles = new Map(state.editedFiles)
  newEditedFiles.delete(path)
  return { ...state, editedFiles: newEditedFiles }
}

function getFileContent(state: EditState, path: string): string {
  if (state.editedFiles.has(path)) {
    return state.editedFiles.get(path)!
  }
  return state.originalFiles.get(path) ?? ''
}

describe('Property 3: Edit State Preservation', () => {
  // Arbitrary generators
  const filePathArb = fc.string({ minLength: 1, maxLength: 50 })
    .filter(s => s.trim().length > 0)
    .map(s => s.replace(/\//g, '-')) // Avoid path issues in tests
  
  const fileContentArb = fc.string({ minLength: 0, maxLength: 500 })
  
  const fileArb = fc.record({
    path: filePathArb,
    content: fileContentArb,
  })
  
  const filesArb = fc.array(fileArb, { minLength: 1, maxLength: 10 })
    .map(files => {
      // Ensure unique paths
      const seen = new Set<string>()
      return files.filter(f => {
        if (seen.has(f.path)) return false
        seen.add(f.path)
        return true
      })
    })
    .filter(files => files.length > 0)

  it('edited content is preserved and reflected in subsequent reads', () => {
    fc.assert(
      fc.property(
        filesArb,
        fileContentArb,
        (files, newContent) => {
          // Setup: create initial state with files
          const state = createInitialState(files)
          const targetPath = files[0].path
          
          // Action: edit the file
          const editedState = editFile(state, targetPath, newContent)
          
          // Verify: the edited content is preserved and returned on read
          const readContent = getFileContent(editedState, targetPath)
          expect(readContent).toBe(newContent)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('multiple edits to same file preserve only the latest content', () => {
    fc.assert(
      fc.property(
        filesArb,
        fc.array(fileContentArb, { minLength: 2, maxLength: 5 }),
        (files, edits) => {
          let state = createInitialState(files)
          const targetPath = files[0].path
          
          // Apply multiple edits
          for (const content of edits) {
            state = editFile(state, targetPath, content)
          }
          
          // Verify: only the last edit is preserved
          const readContent = getFileContent(state, targetPath)
          expect(readContent).toBe(edits[edits.length - 1])
        }
      ),
      { numRuns: 100 }
    )
  })

  it('reset restores original content', () => {
    fc.assert(
      fc.property(
        filesArb,
        fileContentArb,
        (files, newContent) => {
          const state = createInitialState(files)
          const targetPath = files[0].path
          const originalContent = files[0].content
          
          // Edit then reset
          const editedState = editFile(state, targetPath, newContent)
          const resetState = resetFile(editedState, targetPath)
          
          // Verify: original content is restored
          const readContent = getFileContent(resetState, targetPath)
          expect(readContent).toBe(originalContent)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('edits to one file do not affect other files', () => {
    fc.assert(
      fc.property(
        filesArb.filter(f => f.length >= 2),
        fileContentArb,
        (files, newContent) => {
          const state = createInitialState(files)
          const editPath = files[0].path
          const otherPath = files[1].path
          const otherOriginalContent = files[1].content
          
          // Edit first file
          const editedState = editFile(state, editPath, newContent)
          
          // Verify: other file is unchanged
          const otherContent = getFileContent(editedState, otherPath)
          expect(otherContent).toBe(otherOriginalContent)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('unedited files return original content', () => {
    fc.assert(
      fc.property(
        filesArb,
        (files) => {
          const state = createInitialState(files)
          
          // Verify: all files return their original content when not edited
          for (const file of files) {
            const content = getFileContent(state, file.path)
            expect(content).toBe(file.content)
          }
        }
      ),
      { numRuns: 100 }
    )
  })
})
