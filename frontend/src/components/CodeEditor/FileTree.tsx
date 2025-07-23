import React, { useState, useCallback, useEffect } from 'react';
import { 
  ChevronRight, 
  ChevronDown, 
  File, 
  Folder, 
  FolderOpen,
  Search,
  Plus,
  Trash2,
  Edit,
  Download
} from 'lucide-react';
import './FileTree.scss';

export interface FileTreeNode {
  id: string;
  name: string;
  type: 'file' | 'directory';
  path: string;
  children?: FileTreeNode[];
  size?: number;
  lastModified?: Date;
  language?: string;
  content?: string;
}

export interface FileTreeProps {
  files: FileTreeNode[];
  onFileSelect: (file: FileTreeNode) => void;
  onFileCreate?: (parentPath: string, name: string, type: 'file' | 'directory') => void;
  onFileDelete?: (file: FileTreeNode) => void;
  onFileRename?: (file: FileTreeNode, newName: string) => void;
  onFileDownload?: (file: FileTreeNode) => void;
  selectedFile?: string;
  expandedNodes?: Set<string>;
  onExpandedChange?: (expanded: Set<string>) => void;
  searchQuery?: string;
  onSearchChange?: (query: string) => void;
  showSearch?: boolean;
  allowEdit?: boolean;
  className?: string;
}

const FileTree: React.FC<FileTreeProps> = ({
  files,
  onFileSelect,
  onFileCreate,
  onFileDelete,
  onFileRename,
  onFileDownload,
  selectedFile,
  expandedNodes = new Set(),
  onExpandedChange,
  searchQuery = '',
  onSearchChange,
  showSearch = true,
  allowEdit = false,
  className = ''
}) => {
  const [localExpanded, setLocalExpanded] = useState<Set<string>>(expandedNodes);
  const [editingNode, setEditingNode] = useState<string | null>(null);
  const [newNodeName, setNewNodeName] = useState('');
  const [contextMenu, setContextMenu] = useState<{
    x: number;
    y: number;
    node: FileTreeNode;
  } | null>(null);

  const expanded = onExpandedChange ? expandedNodes : localExpanded;
  const setExpanded = onExpandedChange || setLocalExpanded;

  const toggleExpanded = useCallback((nodeId: string) => {
    const newExpanded = new Set(expanded);
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId);
    } else {
      newExpanded.add(nodeId);
    }
    setExpanded(newExpanded);
  }, [expanded, setExpanded]);

  const getFileIcon = useCallback((file: FileTreeNode) => {
    if (file.type === 'directory') {
      return expanded.has(file.id) ? <FolderOpen size={16} /> : <Folder size={16} />;
    }

    // Return specific icons based on file extension
    const extension = file.name.split('.').pop()?.toLowerCase();
    switch (extension) {
      case 'go':
        return <File size={16} className="file-icon--go" />;
      case 'json':
        return <File size={16} className="file-icon--json" />;
      case 'md':
        return <File size={16} className="file-icon--markdown" />;
      case 'txt':
        return <File size={16} className="file-icon--text" />;
      case 'yml':
      case 'yaml':
        return <File size={16} className="file-icon--yaml" />;
      case 'dockerfile':
        return <File size={16} className="file-icon--docker" />;
      default:
        return <File size={16} />;
    }
  }, [expanded]);

  const formatFileSize = useCallback((bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }, []);

  const filterFiles = useCallback((nodes: FileTreeNode[], query: string): FileTreeNode[] => {
    if (!query) return nodes;

    const filterNode = (node: FileTreeNode): FileTreeNode | null => {
      const matches = node.name.toLowerCase().includes(query.toLowerCase());
      
      if (node.type === 'directory' && node.children) {
        const filteredChildren = node.children
          .map(child => filterNode(child))
          .filter(Boolean) as FileTreeNode[];
        
        if (filteredChildren.length > 0 || matches) {
          return {
            ...node,
            children: filteredChildren
          };
        }
      } else if (matches) {
        return node;
      }
      
      return null;
    };

    return nodes.map(node => filterNode(node)).filter(Boolean) as FileTreeNode[];
  }, []);

  const handleContextMenu = useCallback((e: React.MouseEvent, node: FileTreeNode) => {
    if (!allowEdit) return;
    
    e.preventDefault();
    setContextMenu({
      x: e.clientX,
      y: e.clientY,
      node
    });
  }, [allowEdit]);

  const handleContextMenuAction = useCallback((action: string, node: FileTreeNode) => {
    setContextMenu(null);
    
    switch (action) {
      case 'rename':
        setEditingNode(node.id);
        setNewNodeName(node.name);
        break;
      case 'delete':
        if (onFileDelete) {
          onFileDelete(node);
        }
        break;
      case 'download':
        if (onFileDownload) {
          onFileDownload(node);
        }
        break;
      case 'newFile':
        if (onFileCreate) {
          const name = prompt('Enter file name:');
          if (name) {
            onFileCreate(node.path, name, 'file');
          }
        }
        break;
      case 'newFolder':
        if (onFileCreate) {
          const name = prompt('Enter folder name:');
          if (name) {
            onFileCreate(node.path, name, 'directory');
          }
        }
        break;
    }
  }, [onFileDelete, onFileDownload, onFileCreate]);

  const handleRename = useCallback((nodeId: string, newName: string) => {
    if (onFileRename) {
      const findNode = (nodes: FileTreeNode[]): FileTreeNode | null => {
        for (const node of nodes) {
          if (node.id === nodeId) return node;
          if (node.children) {
            const found = findNode(node.children);
            if (found) return found;
          }
        }
        return null;
      };

      const node = findNode(files);
      if (node) {
        onFileRename(node, newName);
      }
    }
    setEditingNode(null);
    setNewNodeName('');
  }, [files, onFileRename]);

  const renderNode = useCallback((node: FileTreeNode, level = 0) => {
    const isExpanded = expanded.has(node.id);
    const isSelected = selectedFile === node.id;
    const isEditing = editingNode === node.id;

    return (
      <div key={node.id} className="file-tree__node">
        <div
          className={`file-tree__item ${isSelected ? 'file-tree__item--selected' : ''}`}
          style={{ paddingLeft: `${level * 20 + 8}px` }}
          onClick={() => {
            if (node.type === 'directory') {
              toggleExpanded(node.id);
            } else {
              onFileSelect(node);
            }
          }}
          onContextMenu={(e) => handleContextMenu(e, node)}
        >
          {node.type === 'directory' && (
            <button
              className="file-tree__expand-button"
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(node.id);
              }}
            >
              {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
            </button>
          )}
          
          <div className="file-tree__icon">
            {getFileIcon(node)}
          </div>
          
          {isEditing ? (
            <input
              className="file-tree__edit-input"
              value={newNodeName}
              onChange={(e) => setNewNodeName(e.target.value)}
              onBlur={() => handleRename(node.id, newNodeName)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleRename(node.id, newNodeName);
                } else if (e.key === 'Escape') {
                  setEditingNode(null);
                  setNewNodeName('');
                }
              }}
              autoFocus
            />
          ) : (
            <>
              <span className="file-tree__name">{node.name}</span>
              {node.type === 'file' && node.size && (
                <span className="file-tree__size">{formatFileSize(node.size)}</span>
              )}
            </>
          )}
        </div>
        
        {node.type === 'directory' && isExpanded && node.children && (
          <div className="file-tree__children">
            {node.children.map(child => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    );
  }, [
    expanded, 
    selectedFile, 
    editingNode, 
    newNodeName, 
    toggleExpanded, 
    onFileSelect, 
    handleContextMenu, 
    getFileIcon, 
    formatFileSize, 
    handleRename
  ]);

  // Close context menu on outside click
  useEffect(() => {
    const handleClickOutside = () => setContextMenu(null);
    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, []);

  const filteredFiles = filterFiles(files, searchQuery);

  return (
    <div className={`file-tree ${className}`}>
      {showSearch && (
        <div className="file-tree__search">
          <div className="file-tree__search-input">
            <Search size={16} />
            <input
              type="text"
              placeholder="Search files..."
              value={searchQuery}
              onChange={(e) => onSearchChange?.(e.target.value)}
            />
          </div>
        </div>
      )}
      
      <div className="file-tree__content">
        {filteredFiles.length > 0 ? (
          filteredFiles.map(node => renderNode(node))
        ) : (
          <div className="file-tree__empty">
            {searchQuery ? 'No files match your search' : 'No files available'}
          </div>
        )}
      </div>

      {contextMenu && (
        <div
          className="file-tree__context-menu"
          style={{ left: contextMenu.x, top: contextMenu.y }}
        >
          {contextMenu.node.type === 'directory' && (
            <>
              <button onClick={() => handleContextMenuAction('newFile', contextMenu.node)}>
                <Plus size={14} />
                New File
              </button>
              <button onClick={() => handleContextMenuAction('newFolder', contextMenu.node)}>
                <Plus size={14} />
                New Folder
              </button>
              <div className="file-tree__context-menu-separator" />
            </>
          )}
          
          <button onClick={() => handleContextMenuAction('rename', contextMenu.node)}>
            <Edit size={14} />
            Rename
          </button>
          
          <button onClick={() => handleContextMenuAction('download', contextMenu.node)}>
            <Download size={14} />
            Download
          </button>
          
          <div className="file-tree__context-menu-separator" />
          
          <button 
            className="file-tree__context-menu-danger"
            onClick={() => handleContextMenuAction('delete', contextMenu.node)}
          >
            <Trash2 size={14} />
            Delete
          </button>
        </div>
      )}
    </div>
  );
};

export default FileTree;