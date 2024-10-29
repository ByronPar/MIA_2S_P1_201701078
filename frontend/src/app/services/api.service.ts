import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private baseUrl = 'http://localhost:8088';

  constructor(
    private httpClient: HttpClient
  ) { }

  postEntrada(entrada: string): Observable<string> {
    //return this.httpClient.post("http://107.20.78.66:8088/comando", entrada, {
      return this.httpClient.post("http://localhost:8088/comando", entrada, {
      headers: { 'Content-Type': 'text/plain' },
      responseType: 'text' // Establecer el tipo de respuesta como texto
    });
  }

  getDatos(path: string ): Observable<any> {
    // Convertir el objeto a JSON antes de enviarlo
    return this.httpClient.post("http://localhost:8088/archivos", path, {
      headers: { 'Content-Type': 'text/plain' },
      responseType: 'json' // Establecer el tipo de respuesta como JSON
    });
  }

  getFileContent(path: string): Observable<string> {
    return this.httpClient.get<string>(`${this.baseUrl}/file-content?path=${encodeURIComponent(path)}`, {
      responseType: 'text' as 'json' // Cast para evitar problemas de tipo
    });
  }


}
