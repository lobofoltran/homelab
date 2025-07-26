interface CameraViewerProps {
  cameraId: string;
  name?: string;
}

export function CameraViewer({ cameraId, name }: CameraViewerProps) {
  const src = `http://localhost:8080/cameras/${cameraId}?t=${Date.now()}`;

  return (
    <div className="bg-muted border border-border rounded-2xl shadow-lg overflow-hidden">
      <div className="px-4 py-2 border-b border-border bg-muted text-muted-foreground text-sm font-medium">
        {name || `Camera ${cameraId}`}
      </div>
      <img
        src={src}
        alt={`Camera ${cameraId}`}
        className="w-full h-[800px] object-cover bg-black"
      />
    </div>
  );
}
