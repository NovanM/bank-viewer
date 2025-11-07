
"use client";

import React, { useRef, useState } from "react";

interface Props {
  onUpload: (file: File) => void;
  isPending: boolean;
}

export function FileUploader({ onUpload, isPending }: Props) {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleFileSelection = (file: File | null) => {
    setSelectedFile(file);
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (selectedFile) {
      onUpload(selectedFile);
    }
  };

  const handleBrowseClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null;
    handleFileSelection(file);
  };

  const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setIsDragging(false);
    const file = event.dataTransfer.files?.[0] ?? null;
    if (file) {
      handleFileSelection(file);
    }
  };

  return (
    <section className="component-section">
      <h2>Upload Statement</h2>
      <p className="component-description">Drop a CSV export or browse your files to refresh the latest statement.</p>
      <form onSubmit={handleSubmit} className="uploader-form">
        <input
          ref={fileInputRef}
          type="file"
          accept=".csv"
          onChange={handleFileChange}
          disabled={isPending}
          hidden
        />
        <div
          className={`dropzone${isDragging ? " dropzone--active" : ""}${selectedFile ? " dropzone--ready" : ""}`}
          onClick={handleBrowseClick}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <p className="dropzone-title">Drag and drop your CSV file here</p>
          <p className="dropzone-subtitle">or click to browse from your device</p>
          {selectedFile && <p className="file-name">Selected: {selectedFile.name}</p>}
        </div>
        <button
          type="submit"
          className="button button--primary"
          disabled={!selectedFile || isPending}
        >
          {isPending ? "Uploading..." : "Upload CSV"}
        </button>
      </form>
    </section>
  );
}