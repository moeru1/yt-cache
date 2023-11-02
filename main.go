package main

import ( 
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "sync"
    //"syscall"
)

func url_to_filename(url string) string {
    return url
}

func download_video(url string) (string, error) {
    //var stat syscall.Statfs_t
    //wd, err_os := os.Getwd()
    //if err_os != nil {
    //    return "", err_os 
    //}

    //syscall.Statfs(wd, &stat)

    //maxsize := (stat.Bavail * uint64(stat.Bsize)) / 2

    //err_mkdir := exec.Command("bash", "-c", "mkdir -p /tmp/dv").Run()

    //if err_mkdir != nil {
    //    return "", err_mkdir
    //}

    //yt-dlp -ic -o "/tmp/dv/dvout-%(playlist-index)s-%(id)s.%(ext)s" "$link" -N8 -f "bv*[filesize<${maxsize}]+ba / b[filesize<${maxsize}] / w" --add-metadata
    cmd := exec.Command("./script.sh", url)

    output, err := cmd.Output()
    filepath := string(output)

    fmt.Println("filepath: " + filepath)

    if err != nil {
        return "", err 
    }

    return filepath, nil 
}

func main() {
    cache := Cache{v: make(map[string]string)}
    in := bufio.NewScanner(os.Stdin)
    var wg sync.WaitGroup
    for in.Scan() {
        url := in.Text()
		_, found := cache.LoadOrStore(url, "")
        if !found {
            wg.Add(1)
            go func() { 
                defer wg.Done()
                filepath, err := download_video(url)
                if err != nil {
                    fmt.Errorf("Error downloading %s: %v", url, err)
                }
                fmt.Println("Downloaded: ", url, " to ", filepath)
            }()
        } else {
            fmt.Println("Already downloaded: ", url)
        }
    }
    wg.Wait()
}
