export const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }
export const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('de-DE', {
      day: '2-digit',
      month: 'long',
      year: 'numeric'
    })
  }

  interface Icons {
    [key: string]: string;
  }
  
  const fileIcons: Icons = {
    zip: "folder_zip",
    rar: "folder_zip",
    '7z': "folder_zip",
    htm: "html",
    html: "html",
    js: "javascript",
    json: "data_object",
    md: "article",
    pdf: "picture_as_pdf",
    png: "image",
    jpg: "image",
    jpeg: "image",
    webp: "image",
    mp4: "movie",
    mkv: "movie",
    avi: "movie",
    wmv: "movie",
    mov: "movie",
    txt: "description",
    xls: "table_chart",
    other: "insert_drive_file"
  };
  
  export const getFileIcon = (fileName: string) => {
    const ext = fileName.split('.').pop()
    if (!ext) return fileIcons['other']
    return fileIcons[ext] ?? fileIcons['other']
  }