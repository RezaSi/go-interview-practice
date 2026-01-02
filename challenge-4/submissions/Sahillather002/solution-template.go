package main

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func BFS(graph map[int][]int, start int)[]int{
    visited:=make(map[int]bool)
    queue:=[]int{start}
    visited[start]=true
    order:=[]int{}
    
    for len(queue)>0{
        u:=queue[0]
        queue=queue[1:]
        order=append(order,u)
        
        for _,v:=range graph[u]{
            if !visited[v]{
                visited[v]=true
                queue=append(queue,v)
            }
        }
    }
    return order
}
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	type result struct{
	    query int
	    order []int
	}
	jobs := make(chan int, len(queries))
	results := make(chan result, len(queries))
	
	if numWorkers<=0{
	    return map[int][]int{}
	}
	
	//worker
	worker:=func(){
	    for q:=range jobs{
	        order:=BFS(graph,q)
	        results<-result{query:q,order:order}
	    }
	}
	
	//launch workers
	for i:=0;i<numWorkers;i++{
	    go worker()
	}
	
	//send jobs
	for _,q:=range queries{
	    jobs<-q
	}
	close(jobs)
	
	//collect results
	output:=make(map[int][]int,len(queries))
	for i:=0;i<len(queries);i++{
	    r:=<-results
	    output[r.query]=r.order
	}
	return output
}

func main() {
	// You can insert optional local tests here if desired.
}
