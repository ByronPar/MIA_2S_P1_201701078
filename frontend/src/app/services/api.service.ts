import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(
    private httpClient: HttpClient
  ) { }

  postEntrada(entrada: string) {
    // Cambiar la URL por la de la API en AWS
    return this.httpClient.post("http://107.20.78.66:8080/comando", { Cmd: entrada });
  }
}
