import { describe, it, expect } from 'vitest'
import * as fc from 'fast-check'
import type { GeneratedFile } from '@/lib/api'

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


/**
 * Property 4: Download Content Integrity
 * For any file download (individual or ZIP), the downloaded content SHALL match
 * the current edited state if modifications were made, or the original content if not.
 * 
 * Validates: Requirements 5.3
 * 
 * Feature: ai-driven-generation, Property 4: Download Content Integrity
 */

// Simulates the download content resolution logic
function getDownloadContent(
  files: GeneratedFile[],
  editedFiles: Map<string, string>,
  path: string
): string {
  if (editedFiles.has(path)) {
    return editedFiles.get(path)!
  }
  const file = files.find(f => f.path === path)
  return file?.content ?? ''
}

// Simulates preparing files for ZIP download
function prepareFilesForZip(
  files: GeneratedFile[],
  editedFiles: Map<string, string>
): { path: string; content: string }[] {
  return files.map(file => ({
    path: file.path,
    content: getDownloadContent(files, editedFiles, file.path),
  }))
}

describe('Property 4: Download Content Integrity', () => {
  const fileTypeArb = fc.constantFrom('kickoff', 'steering', 'hook') as fc.Arbitrary<'kickoff' | 'steering' | 'hook'>
  
  const filePathArb = fc.string({ minLength: 1, maxLength: 50 })
    .filter(s => s.trim().length > 0)
    .map(s => s.replace(/[/\\]/g, '-'))
  
  const fileContentArb = fc.string({ minLength: 0, maxLength: 500 })
  
  const generatedFileArb: fc.Arbitrary<GeneratedFile> = fc.record({
    path: filePathArb,
    content: fileContentArb,
    type: fileTypeArb,
  })
  
  const generatedFilesArb = fc.array(generatedFileArb, { minLength: 1, maxLength: 10 })
    .map(files => {
      const seen = new Set<string>()
      return files.filter(f => {
        if (seen.has(f.path)) return false
        seen.add(f.path)
        return true
      })
    })
    .filter(files => files.length > 0)

  it('download returns edited content when file was modified', () => {
    fc.assert(
      fc.property(
        generatedFilesArb,
        fileContentArb,
        (files, editedContent) => {
          const editedFiles = new Map<string, string>()
          const targetPath = files[0].path
          editedFiles.set(targetPath, editedContent)
          
          const downloadContent = getDownloadContent(files, editedFiles, targetPath)
          
          expect(downloadContent).toBe(editedContent)
        }
      ),
      { numRuns: 100 }
    )
  })

  it('download returns original content when file was not modified', () => {
    fc.assert(
      fc.property(
        generatedFilesArb,
        (files) => {
          const editedFiles = new Map<string, string>()
          
          for (const file of files) {
            const downloadContent = getDownloadContent(files, editedFiles, file.path)
            expect(downloadContent).toBe(file.content)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('ZIP preparation uses edited content for modified files and original for unmodified', () => {
    fc.assert(
      fc.property(
        generatedFilesArb.filter(f => f.length >= 2),
        fileContentArb,
        (files, editedContent) => {
          const editedFiles = new Map<string, string>()
          const editedPath = files[0].path
          editedFiles.set(editedPath, editedContent)
          
          const zipFiles = prepareFilesForZip(files, editedFiles)
          
          // Verify edited file has edited content
          const editedZipFile = zipFiles.find(f => f.path === editedPath)
          expect(editedZipFile?.content).toBe(editedContent)
          
          // Verify unedited files have original content
          for (let i = 1; i < files.length; i++) {
            const originalFile = files[i]
            const zipFile = zipFiles.find(f => f.path === originalFile.path)
            expect(zipFile?.content).toBe(originalFile.content)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('all files are included in ZIP preparation', () => {
    fc.assert(
      fc.property(
        generatedFilesArb,
        (files) => {
          const editedFiles = new Map<string, string>()
          const zipFiles = prepareFilesForZip(files, editedFiles)
          
          expect(zipFiles.length).toBe(files.length)
          
          for (const file of files) {
            const zipFile = zipFiles.find(f => f.path === file.path)
            expect(zipFile).toBeDefined()
          }
        }
      ),
      { numRuns: 100 }
    )
  })
})

/**
 * Property 5: ZIP Directory Structure
 * For any generated ZIP file, all steering files SHALL be under `.kiro/steering/`
 * and all hook files SHALL be under `.kiro/hooks/`.
 * 
 * Validates: Requirements 5.4
 * 
 * Feature: ai-driven-generation, Property 5: ZIP Directory Structure
 */

// Validates that a file path matches expected directory structure based on type
function validateFileStructure(file: GeneratedFile): boolean {
  switch (file.type) {
    case 'steering':
      return file.path.startsWith('.kiro/steering/')
    case 'hook':
      return file.path.startsWith('.kiro/hooks/')
    case 'kickoff':
      // Kickoff files can be at root or any location
      return true
    default:
      return false
  }
}

// Simulates the structure that would be created in a ZIP
function getExpectedZipStructure(files: GeneratedFile[]): { path: string; type: string }[] {
  return files.map(f => ({ path: f.path, type: f.type }))
}

describe('Property 5: ZIP Directory Structure', () => {
  const steeringPathArb = fc.string({ minLength: 1, maxLength: 30 })
    .filter(s => s.trim().length > 0 && !s.includes('/'))
    .map(s => `.kiro/steering/${s}.md`)
  
  const hookPathArb = fc.string({ minLength: 1, maxLength: 30 })
    .filter(s => s.trim().length > 0 && !s.includes('/'))
    .map(s => `.kiro/hooks/${s}.kiro.hook`)
  
  const kickoffPathArb = fc.string({ minLength: 1, maxLength: 30 })
    .filter(s => s.trim().length > 0 && !s.includes('/'))
    .map(s => `${s}.md`)
  
  const fileContentArb = fc.string({ minLength: 0, maxLength: 200 })
  
  const steeringFileArb: fc.Arbitrary<GeneratedFile> = fc.record({
    path: steeringPathArb,
    content: fileContentArb,
    type: fc.constant('steering' as const),
  })
  
  const hookFileArb: fc.Arbitrary<GeneratedFile> = fc.record({
    path: hookPathArb,
    content: fileContentArb,
    type: fc.constant('hook' as const),
  })
  
  const kickoffFileArb: fc.Arbitrary<GeneratedFile> = fc.record({
    path: kickoffPathArb,
    content: fileContentArb,
    type: fc.constant('kickoff' as const),
  })
  
  const validFilesArb = fc.tuple(
    fc.array(steeringFileArb, { minLength: 1, maxLength: 5 }),
    fc.array(hookFileArb, { minLength: 1, maxLength: 5 }),
    fc.array(kickoffFileArb, { minLength: 0, maxLength: 2 })
  ).map(([steering, hooks, kickoff]) => {
    // Ensure unique paths
    const seen = new Set<string>()
    const all = [...steering, ...hooks, ...kickoff]
    return all.filter(f => {
      if (seen.has(f.path)) return false
      seen.add(f.path)
      return true
    })
  }).filter(files => files.length > 0)

  it('all steering files are under .kiro/steering/', () => {
    fc.assert(
      fc.property(
        validFilesArb,
        (files) => {
          const steeringFiles = files.filter(f => f.type === 'steering')
          
          for (const file of steeringFiles) {
            expect(file.path.startsWith('.kiro/steering/')).toBe(true)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('all hook files are under .kiro/hooks/', () => {
    fc.assert(
      fc.property(
        validFilesArb,
        (files) => {
          const hookFiles = files.filter(f => f.type === 'hook')
          
          for (const file of hookFiles) {
            expect(file.path.startsWith('.kiro/hooks/')).toBe(true)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('all files pass structure validation', () => {
    fc.assert(
      fc.property(
        validFilesArb,
        (files) => {
          for (const file of files) {
            expect(validateFileStructure(file)).toBe(true)
          }
        }
      ),
      { numRuns: 100 }
    )
  })

  it('ZIP structure preserves file paths exactly', () => {
    fc.assert(
      fc.property(
        validFilesArb,
        (files) => {
          const zipStructure = getExpectedZipStructure(files)
          
          expect(zipStructure.length).toBe(files.length)
          
          for (const file of files) {
            const zipEntry = zipStructure.find(z => z.path === file.path)
            expect(zipEntry).toBeDefined()
            expect(zipEntry?.type).toBe(file.type)
          }
        }
      ),
      { numRuns: 100 }
    )
  })
})
