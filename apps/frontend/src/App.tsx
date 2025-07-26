import { CameraViewer } from "@/components/camera-viewer";

export default function App() {
  return (
    <div className="min-h-svh bg-background text-foreground flex flex-col items-center justify-center gap-8 p-8 dark">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 w-full max-w-8xl">
        <CameraViewer cameraId="0" name="Sacada" />
        <CameraViewer cameraId="1" name="Frente Loja" />
      </div>
    </div>
  );
}
